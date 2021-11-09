package job_mgr

import (
	"context"
	"cronTab/common"
	"cronTab/common/cron_job"
	"cronTab/worker/etcd_ops"
	terrors "github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	jobMgr *JobMgr
)

// InitJobMgr 初始化 JobMgr
func InitJobMgr() (*JobMgr, error) {
	var err error

	// 获取 kv 和 lease
	jobMgr = &JobMgr{
		kv:      etcd_ops.GetKv(),
		lease:   etcd_ops.GetLease(),
		watcher: etcd_ops.GetWatcher(),
	}

	// 启动执行器
	if err = InitExecutor(); err != nil {
		return nil, err
	}

	// 启动调度
	if err = InitScheduler(); err != nil {
		return nil, err
	}

	return jobMgr, nil
}

// GetJobRevision 监听任务
func (jobMgr *JobMgr) GetJobRevision() (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 1. Get /cron/jobs/ 目录下的所有任务, 并且获得当前集群的 Revision
	getResp, err := jobMgr.kv.Get(ctx, common.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return 0, terrors.Wrap(err, "get jobs info from etcd failed")
	}

	// 当前的任务
	for _, kv := range getResp.Kvs {
		if job, err := cron_job.JobUnmarshal(kv.Value); err != nil {
			// 值非法
			continue
		} else {
			jobEvent := cron_job.BuildJobEvent(common.JobEventSave, job)
			// 任务同步给 scheduler(调度协程)
			GScheduler.PushJobEvent(jobEvent)
		}
	}

	// 从 GET 时刻的后续版本开始监听变化
	watchStartRevision := getResp.Header.Revision + 1

	return watchStartRevision, nil
}

// WatchJob 监听任务变化函数
func (jobMgr *JobMgr) WatchJob(watchStartRevision int64, stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 监听 /cron/jobs/ 目录的后续变化
	watchChan := jobMgr.watcher.Watch(ctx, common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

	go func() {
		<-stop
		cancel()
	}()

	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			var jobEvent *cron_job.JobEvent
			switch ev.Type {
			case mvccpb.PUT: // 任务保存事件
				job, err := cron_job.JobUnmarshal(ev.Kv.Value)
				if err != nil {
					continue
				}
				// 构造一个更新 event 事件
				jobEvent = cron_job.BuildJobEvent(common.JobEventSave, job)
			case mvccpb.DELETE: // 任务删除事件
				jobName := cron_job.ExtractJobName(string(ev.Kv.Key))
				// 构造一个删除 event 事件
				jobEvent = cron_job.BuildJobEvent(common.JobEventDelete, &cron_job.Job{Name: jobName})
			}
			// 推送事件给 scheduler(调度协程)
			GScheduler.PushJobEvent(jobEvent)
		}
	}

	return nil
}

// WatchKill 监听强杀变化函数
func (jobMgr *JobMgr) WatchKill(stop <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 监听 /cron/kill/ 目录的变化
	watchChan := jobMgr.watcher.Watch(ctx, common.JobKillDir, clientv3.WithPrefix())

	go func() {
		<-stop
		cancel()
	}()

	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			switch ev.Type {
			case mvccpb.PUT: // 杀死任务事件
				jobName := cron_job.ExtractKillName(string(ev.Kv.Key))
				jobEvent := cron_job.BuildJobEvent(common.JobEventKill, &cron_job.Job{Name: jobName})
				// 推送事件给 scheduler(调度协程)
				GScheduler.PushJobEvent(jobEvent)
			case mvccpb.DELETE: // kill 自动过期
			}
		}
	}

	return nil
}
