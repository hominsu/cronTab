package etcdOps

import (
	"context"
	"cronTab/common"
	"encoding/json"
	"errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

// SaveJob 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (*common.Job, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 任务保存到: /cron/jobs/任务名 -> json

	// etcd 的保存 key
	jobKey := common.JobSaveDir + job.Name

	// 任务信息
	jobValue, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	// 保存到 etcd
	putResp, err := jobMgr.kv.Put(ctx, jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	// 如果是更新，返回旧值
	if putResp.PrevKv != nil {
		oldJob := &common.Job{}
		// 对旧值进行反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, oldJob); err != nil {
			// 如果旧值非法，不影响，因为 Put 成功了
			return nil, nil
		}
		return oldJob, nil
	} else {
		return nil, nil
	}
}

// DeleteJob 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (*common.Job, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// etcd 的保存 key
	jobKey := common.JobSaveDir + name

	// 删除任务
	delResp, err := jobMgr.kv.Delete(ctx, jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	// 如果删除的 key 存在
	if delResp.PrevKvs != nil {
		oldJob := &common.Job{}
		// 对旧值进行反序列化
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, oldJob); err != nil {
			// 如果旧值非法，不影响，因为 Del 成功了
			return nil, nil
		}
		return oldJob, nil
	} else {
		return nil, errors.New("the deleted key does not exist")
	}
}

// ListJobs 列举任务
func (jobMgr *JobMgr) ListJobs() ([]*common.Job, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// etcd 任务目录
	dirKey := common.JobSaveDir

	// 获取目录下所有任务信息
	getResp, err := jobMgr.kv.Get(ctx, dirKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var jobs []*common.Job
	// 遍历所有任务，反序列化
	for _, kv := range getResp.Kvs {
		job := &common.Job{}
		// 对旧值进行反序列化, 这里容忍了错误
		if err = json.Unmarshal(kv.Value, job); err != nil {
			// 值非法
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// KillJob 杀死任务
func (jobMgr *JobMgr) KillJob(name string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 通知 worker 杀死对应的任务
	killKey := common.JobKillDir + name

	// 让 worker 监听一次 put 操作, 创建一个租约让其稍后过期
	leaseGrantResp, err := jobMgr.lease.Grant(ctx, 1)
	if err != nil {
		return err
	}

	// 租约 id
	leaseId := leaseGrantResp.ID

	// 设置 kill 标记
	_, err = jobMgr.kv.Put(ctx, killKey, "", clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}

	return nil
}
