package etcdOps

import (
	"context"
	"cronTab/common"
	"cronTab/worker/scheduler"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func (jobMgr *JobMgr) WatchJob() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 1. Get /cron/jobs/ 目录下的所有任务, 并且获得当前集群的 Revision
	getResp, err := jobMgr.kv.Get(ctx, common.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 当前的任务
	for _, kv := range getResp.Kvs {
		if job, err := common.JobUnmarshal(kv.Value); err != nil {
			// 值非法
			continue
		} else {
			jobEvent := common.BuildJobEvent(common.JobEventSave, job)
			// 任务同步给 scheduler(调度协程)
			scheduler.GScheduler.PushJobEvent(jobEvent)
		}
	}

	// 从 GET 时刻的后续版本开始监听变化
	watchStartRevision := getResp.Header.Revision + 1

	// 监听协程
	go watch(jobMgr, watchStartRevision)

	return nil
}

func watch(jobMgr *JobMgr, watchStartRevision int64) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// 监听 /cron/jobs/ 目录的后续变化
	watchChan := jobMgr.watcher.Watch(ctx, common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

	for watchResp := range watchChan {
		for _, ev := range watchResp.Events {
			var jobEvent *common.JobEvent
			switch ev.Type {
			case mvccpb.PUT: // 任务保存事件
				job, err := common.JobUnmarshal(ev.Kv.Value)
				if err != nil {
					continue
				}
				// 构造一个更新 event 事件
				jobEvent = common.BuildJobEvent(common.JobEventSave, job)
			case mvccpb.DELETE: // 任务删除事件
				jobName := common.ExtractJobName(string(ev.Kv.Key))
				// 构造一个删除 event 事件
				jobEvent = common.BuildJobEvent(common.JobEventDelete, &common.Job{Name: jobName})
			}
			// 推送事件给 scheduler(调度协程)
			scheduler.GScheduler.PushJobEvent(jobEvent)
		}
	}
}
