package etcd_ops

import (
	"context"
	"cronTab/worker/config"
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
	EtcdCli *etcdCli
)

// InitEtcdConn 初始化 etcd 连接
func InitEtcdConn() error {
	var cli *clientv3.Client
	var err error

	ctx, cancel := context.WithTimeout(context.TODO(), 1000*time.Millisecond)
	defer cancel()

	conf := clientv3.Config{
		Endpoints:   config.GConfig.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GConfig.EtcdDailTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if cli, err = clientv3.New(conf); err != nil {
		return terrors.Wrap(err, "create etcd connection failed")
	}

	EtcdCli = &etcdCli{
		cli:     cli,
		kv:      clientv3.NewKV(cli),
		lease:   clientv3.NewLease(cli),
		watcher: clientv3.NewWatcher(cli),
	}

	// 测试 etcd 连接
	if _, err = EtcdCli.kv.Get(ctx, "/cron"); err != nil {
		return terrors.Wrap(err, "test etcd connection failed")
	}

	return nil
}

// CloseEtcdConn 关闭 etcd 连接
func CloseEtcdConn() error {
	if err := EtcdCli.cli.Close(); err != nil {
		return terrors.Wrap(err, "disconnect etcd failed")
	}
	return nil
}

// GetKv 返回 etcd 的 kv
func (e etcdCli) GetKv() clientv3.KV {
	return e.kv
}

// GetLease 返回 etcd 的 lease
func (e etcdCli) GetLease() clientv3.Lease {
	return e.lease
}

// GetWatcher 返回 etcd 的 watcher
func (e etcdCli) GetWatcher() clientv3.Watcher {
	return e.watcher
}
