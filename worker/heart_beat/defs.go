package heart_beat

import (
	"context"
	"cronTab/common"
	"cronTab/worker/etcd_ops"
	"go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

type HeartBeat struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	leaseId    clientv3.LeaseID   // 租约 ID
	cancelFunc context.CancelFunc // 终止自动续租
}

// InitHeartBeat 初始化心跳
func InitHeartBeat() *HeartBeat {
	return &HeartBeat{
		kv:    etcd_ops.EtcdCli.GetKv(),
		lease: etcd_ops.EtcdCli.GetLease(),
	}
}

// StartHeartBeat 开始心跳
func (heartBeat *HeartBeat) StartHeartBeat() error {
	// 1. 创建租约
	leaseGrantResp, err := heartBeat.lease.Grant(context.TODO(), 5)
	if err != nil {
		return err
	}

	// 2. 自动续租
	leaseId := leaseGrantResp.ID
	ctx, cancelFunc := context.WithCancel(context.TODO())

	heartBeat.leaseId = leaseId
	heartBeat.cancelFunc = cancelFunc

	keepRespChan, err := heartBeat.lease.KeepAlive(ctx, leaseId)
	if err != nil {
		if err := heartBeat.revokeLease(); err != nil {
			return err
		}
		return err
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
		return err
	}

	// 节点地址
	nodeIpNetKey := common.NodeIpNet + ipNetStr

	if _, err = heartBeat.kv.Put(context.TODO(), nodeIpNetKey, "", clientv3.WithLease(leaseId)); err != nil {
		if err := heartBeat.revokeLease(); err != nil {
			return err
		}
		return err
	}

	return nil
}

// EndHeartBeat 开始心跳
func (heartBeat *HeartBeat) EndHeartBeat() error {
	if err := heartBeat.revokeLease(); err != nil {
		return err
	}
	return nil
}

// revokeLease 关闭续租并释放租约
func (heartBeat *HeartBeat) revokeLease() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	// 关闭自动续租的协程
	heartBeat.cancelFunc()

	// 释放租约
	if _, err := heartBeat.lease.Revoke(ctx, heartBeat.leaseId); err != nil {
		return err
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
