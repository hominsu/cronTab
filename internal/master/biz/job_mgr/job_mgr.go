package job_mgr

import (
	"context"
	"cronTab/internal/master/data/etcd"
	"cronTab/internal/master/service/log_sink"
	"cronTab/internal/pkg/constants"
	"cronTab/internal/pkg/cron_job"
	"errors"
	terrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"strings"
)

type JobMgr struct {
	kv    clientv3.KV
	lease clientv3.Lease
}

var (
	GJobMgr *JobMgr
)

func InitJobMgr() error {
	var err error

	// 获取 kv 和 lease
	GJobMgr = &JobMgr{
		kv:    etcd.Cli.GetKv(),
		lease: etcd.Cli.GetLease(),
	}

	// 初始化日志池
	if err = log_sink.InitLogSink(); err != nil {
		return err
	}

	return nil
}

// SaveJob 保存任务
func (jobMgr *JobMgr) SaveJob(ctx context.Context, job *cron_job.Job) (*cron_job.Job, error) {

	// 任务保存到: /cron/jobs/任务名 -> json

	// etcd 的保存 key
	jobKey := constants.JobSaveDir + job.Name

	// 任务信息
	jobValue, err := job.JobMarshal()
	if err != nil {
		return nil, terrors.Wrap(err, "marshal save job info failed")
	}

	// 保存到 etcd
	putResp, err := jobMgr.kv.Put(ctx, jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return nil, terrors.Wrap(err, "put job info to etcd failed")
	}

	// 如果是更新，返回旧值
	if putResp.PrevKv != nil {
		// 对旧值进行反序列化
		if oldJob, err := cron_job.JobUnmarshal(putResp.PrevKv.Value); err != nil {
			return nil, nil
		} else {
			return oldJob, nil
		}
	} else {
		return nil, nil
	}
}

// DeleteJob 删除任务
func (jobMgr *JobMgr) DeleteJob(ctx context.Context, name string) (*cron_job.Job, error) {
	// etcd 的保存 key
	jobKey := constants.JobSaveDir + name

	// 删除任务
	delResp, err := jobMgr.kv.Delete(ctx, jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, terrors.Wrap(err, "delete job info from etcd failed")
	}

	// 如果删除的 key 存在
	if delResp.PrevKvs != nil {
		// 对旧值进行反序列化
		if oldJob, err := cron_job.JobUnmarshal(delResp.PrevKvs[0].Value); err != nil {
			return nil, terrors.Wrap(err, "unmarshal delete job info failed")
		} else {
			return oldJob, nil
		}
	} else {
		return nil, errors.New("the deleted key does not exist")
	}
}

// ListJobs 列举任务
func (jobMgr *JobMgr) ListJobs(ctx context.Context) ([]*cron_job.Job, error) {
	// etcd 任务目录
	dirKey := constants.JobSaveDir

	// 获取目录下所有任务信息
	getResp, err := jobMgr.kv.Get(ctx, dirKey, clientv3.WithPrefix())
	if err != nil {
		return nil, terrors.Wrap(err, "get jobs info from etcd failed")
	}

	var jobs []*cron_job.Job
	// 遍历所有任务，反序列化
	for _, kv := range getResp.Kvs {
		// 对旧值进行反序列化, 这里容忍了错误
		job, err := cron_job.JobUnmarshal(kv.Value)
		if err != nil {
			// 值非法
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// KillJob 杀死任务
func (jobMgr *JobMgr) KillJob(ctx context.Context, name string) error {
	// 通知 worker 杀死对应的任务
	killKey := constants.JobKillDir + name

	// 让 worker 监听一次 put 操作, 创建一个租约让其稍后过期
	leaseGrantResp, err := jobMgr.lease.Grant(ctx, 1)
	if err != nil {
		return terrors.Wrap(err, "create kill key lease failed")
	}

	// 租约 id
	leaseId := leaseGrantResp.ID

	// 设置 kill 标记
	_, err = jobMgr.kv.Put(ctx, killKey, "", clientv3.WithLease(leaseId))
	if err != nil {
		return terrors.Wrap(err, "put kill key to etcd failed")
	}

	return nil
}

// ListNodes 列出全部节点
func (jobMgr *JobMgr) ListNodes(ctx context.Context) ([]string, error) {
	getResp, err := jobMgr.kv.Get(ctx, constants.NodeIpNet, clientv3.WithPrefix())
	if err != nil {
		return nil, terrors.Wrap(err, "get nodes info from etcd failed")
	}

	var nodes []string
	for _, kv := range getResp.Kvs {
		nodes = append(nodes, strings.TrimPrefix(string(kv.Key), constants.NodeIpNet))
	}

	return nodes, nil
}
