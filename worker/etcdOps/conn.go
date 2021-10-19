package etcdOps

import (
	"cronTab/worker/config"
	"github.com/golang/glog"
	"go.etcd.io/etcd/client/v3"
	"time"
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

	conf := clientv3.Config{
		Endpoints:   config.GConfig.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GConfig.EtcdDailTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if cli, err := clientv3.New(conf); err != nil {
		return err
	} else {
		GJobMgr.cli = cli
	}

	// 获取 kv 和 lease
	GJobMgr.kv = clientv3.NewKV(GJobMgr.cli)
	GJobMgr.lease = clientv3.NewLease(GJobMgr.cli)
	GJobMgr.watcher = clientv3.NewWatcher(GJobMgr.cli)

	return nil
}

func CloseEtcdConn() {
	if err := GJobMgr.cli.Close(); err != nil {
		glog.Fatal(err)
	}
}
