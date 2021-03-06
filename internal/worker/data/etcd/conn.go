package etcd

import (
	"context"
	"cronTab/configs/worker_conf"
	"cronTab/internal/pkg/sync"
	terrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

// etcdCli etcd 连接
type etcdCli struct {
	cli     *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	etcdC *etcdCli
)

// InitEtcdConn 初始化 etcd 连接
func InitEtcdConn(ctx context.Context) error {
	var cli *clientv3.Client
	var err error

	conf := clientv3.Config{
		Endpoints:   worker_conf.GConfig.EtcdEndpoints,                                                             // 集群地址
		DialTimeout: sync.ShrinkDeadLine(ctx, time.Duration(worker_conf.GConfig.EtcdDailTimeout)*time.Millisecond), // 连接超时
	}

	// 建立连接
	if cli, err = clientv3.New(conf); err != nil {
		return terrors.Wrap(err, "create etcd connection failed")
	}

	etcdC = &etcdCli{
		cli:     cli,
		kv:      clientv3.NewKV(cli),
		lease:   clientv3.NewLease(cli),
		watcher: clientv3.NewWatcher(cli),
	}

	// 测试 etcd 连接
	if _, err = etcdC.kv.Get(ctx, "/cron"); err != nil {
		return terrors.Wrap(err, "test etcd connection failed")
	}

	return nil
}

// CloseEtcdConn 关闭 etcd 连接
func CloseEtcdConn() error {
	if err := etcdC.cli.Close(); err != nil {
		return terrors.Wrap(err, "disconnect etcd failed")
	}
	return nil
}

// GetKv 返回 etcd 的 kv
func GetKv() clientv3.KV {
	return etcdC.kv
}

// GetLease 返回 etcd 的 lease
func GetLease() clientv3.Lease {
	return etcdC.lease
}

// GetWatcher 返回 etcd 的 watcher
func GetWatcher() clientv3.Watcher {
	return etcdC.watcher
}
