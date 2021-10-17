package etcdOps

import (
	config2 "cronTab/master/config"
	"github.com/golang/glog"
	"go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	cli   *clientv3.Client
	kv    clientv3.KV
	lease clientv3.Lease
}

var (
	GJobMgr *JobMgr
)

func InitJobMgr() {
	defer glog.Flush()

	var cli *clientv3.Client
	var err error

	config := clientv3.Config{
		Endpoints:   config2.GConfig.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config2.GConfig.EtcdDailTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if cli, err = clientv3.New(config); err != nil {
		glog.Fatal(err)
	}

	// 获取 kv 和 lease
	kv := clientv3.NewKV(cli)
	lease := clientv3.NewLease(cli)

	// 赋值单例
	GJobMgr = &JobMgr{
		cli:   cli,
		kv:    kv,
		lease: lease,
	}
}
