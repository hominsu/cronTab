package etcdOps

import (
	"cronTab/master/config"
	"github.com/golang/glog"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	cli   *clientv3.Client
	kv    clientv3.KV
	lease clientv3.Lease
}

var (
	GJobMgr JobMgr
)

func InitJobMgr() (err error) {

	conf := clientv3.Config{
		Endpoints:   config.GConfig.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GConfig.EtcdDailTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if GJobMgr.cli, err = clientv3.New(conf); err != nil {
		return err
	}

	// 获取 kv 和 lease
	GJobMgr.kv = clientv3.NewKV(GJobMgr.cli)
	GJobMgr.lease = clientv3.NewLease(GJobMgr.cli)

	return nil
}

func CloseEtcdConn() {
	if err := GJobMgr.cli.Close(); err != nil {
		glog.Fatal(err)
	}
}
