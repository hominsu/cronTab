package heart_beat

import (
	"context"
	"cronTab/internal/pkg/constants"
	"cronTab/internal/worker/data/etcd"
	terrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

type heartBeat struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	leaseId    clientv3.LeaseID   // 租约 ID
	cancelFunc context.CancelFunc // 终止自动续租
}

var (
	hb *heartBeat
)

// InitHeartBeat 初始化心跳
func InitHeartBeat(ctx context.Context) error {
	hb = &heartBeat{
		kv:    etcd.GetKv(),
		lease: etcd.GetLease(),
	}
	if err := hb.startHeartBeat(ctx); err != nil {
		return err
	}
	return nil
}

// StopHeartBeat 停止心跳
func StopHeartBeat() error {
	if err := hb.endHeartBeat(); err != nil {
		return err
	}
	return nil
}

// startHeartBeat 开始心跳
func (heartBeat *heartBeat) startHeartBeat(ctx context.Context) error {
	// 1. 创建租约
	leaseGrantResp, err := heartBeat.lease.Grant(ctx, 5)
	if err != nil {
		return terrors.Wrap(err, "create heart beat lease failed")
	}

	// 2. 自动续租
	leaseId := leaseGrantResp.ID
	ctxLease, cancelFunc := context.WithCancel(context.Background())

	heartBeat.leaseId = leaseId
	heartBeat.cancelFunc = cancelFunc

	keepRespChan, err := heartBeat.lease.KeepAlive(ctxLease, leaseId)
	if err != nil {
		if err := heartBeat.revokeLease(); err != nil {
			return err
		}
		return terrors.Wrap(err, "keep heart beat alive failed")
	}

	// 处理续租应答的协程
	go func() {
		var keepResp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					// 租约已经失效
					return
				}
			}
		}
	}()

	// 获取当前节点 ip 地址
	ipNetStr, err := nodeIpNet()
	if err != nil {
		return terrors.Wrap(err, "get current node ip failed")
	}

	// 节点地址
	nodeIpNetKey := constants.NodeIpNet + ipNetStr

	if _, err = heartBeat.kv.Put(ctx, nodeIpNetKey, "", clientv3.WithLease(leaseId)); err != nil {
		if err := heartBeat.revokeLease(); err != nil {
			return err
		}
		return terrors.Wrap(err, "put heart beat to etcd failed")
	}

	return nil
}

// endHeartBeat 开始心跳
func (heartBeat *heartBeat) endHeartBeat() error {
	if err := heartBeat.revokeLease(); err != nil {
		return terrors.Wrap(err, "stop heart beat failed")
	}
	return nil
}

// revokeLease 关闭续租并释放租约
func (heartBeat *heartBeat) revokeLease() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 关闭自动续租的协程
	heartBeat.cancelFunc()

	// 释放租约
	if _, err := heartBeat.lease.Revoke(ctx, heartBeat.leaseId); err != nil {
		return terrors.Wrap(err, "revoke heart beat lease failed")
	}
	return nil
}

// 获取当前节点 ip 地址
func nodeIpNet() (string, error) {
	netInterFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, interFace := range netInterFaces {
		if (interFace.Flags & net.FlagUp) != 0 {
			addrs, _ := interFace.Addrs()
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", nil
}
