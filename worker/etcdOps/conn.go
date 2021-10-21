package etcdOps

import (
	"context"
	"cronTab/worker/config"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var (
	EtcdCli *clientv3.Client
)

// InitEtcdConn 初始化 etcd 连接
func InitEtcdConn() error {
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conf := clientv3.Config{
		Endpoints:   config.GConfig.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GConfig.EtcdDailTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if EtcdCli, err = clientv3.New(conf); err != nil {
		return err
	}

	// 测试 etcd 连接
	if _, err = EtcdCli.KV.Get(ctx, "/cron"); err != nil {
		return err
	}

	return nil
}
