package jobMgr

import (
	"cronTab/worker/etcdOps"
	"github.com/golang/glog"
	"go.etcd.io/etcd/client/v3"
)

type JobMgr struct {
	cli     *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	GJobMgr JobMgr
)

func InitJobMgr() error {

	GJobMgr.cli = etcdOps.EtcdCli

	// 获取 kv 和 lease
	GJobMgr.kv = clientv3.NewKV(etcdOps.EtcdCli)
	GJobMgr.lease = clientv3.NewLease(etcdOps.EtcdCli)
	GJobMgr.watcher = clientv3.NewWatcher(etcdOps.EtcdCli)

	return nil
}

func CloseEtcdConn() {
	if err := GJobMgr.cli.Close(); err != nil {
		glog.Fatal(err)
	}
}
