package etcdOps

import (
	"context"
	"cronTab/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func (jobMgr JobMgr) WatchJob() error {
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
			//TODO: 把任务同步给 scheduler(调度协程)
		}
	}

	// 2. 从该 Revision 向后监听事件
	// 从 GET 时刻的后续版本开始监听变化
	watchStartRevision := getResp.Header.Revision + 1

	// 监听协程
	go func(watchStartRevision int64, saveDir string) {
		ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
		defer cancel()

		// 监听 /cron/jobs/ 目录的后续变化
		watchChan := jobMgr.watcher.Watch(ctx, saveDir, clientv3.WithRev(watchStartRevision))

		for watchResp := range watchChan {
			for _, ev := range watchResp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					// TODO: 反序列 job，推送更新事件给 scheduler(调度协程)
					if job, err := common.JobUnmarshal(ev.Kv.Value); err != nil {
						return
					} else {

					}
				case mvccpb.DELETE:
					// TODO: 反序列 job，推送删除事件给 scheduler(调度协程)
					if job, err := common.JobUnmarshal(ev.Kv.Value); err != nil {
						return
					} else {

					}
				}
			}
		}
	}(watchStartRevision, common.JobSaveDir)

	return nil
}
