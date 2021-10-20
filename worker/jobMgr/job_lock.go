package jobMgr

import (
	"context"
	"cronTab/common"
	"go.etcd.io/etcd/client/v3"
)

// JobLock 分布式锁(TXN 事务)
type JobLock struct {
	kv      clientv3.KV
	lease   clientv3.Lease
	JobName string // 任务名

	cancelFunc context.CancelFunc // 终止自动续租
	leaseId    clientv3.LeaseID   // 租约 ID
	isLocked   bool               // 是否上锁成功
}

// InitJobLock 初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) *JobLock {
	return &JobLock{
		kv:       kv,
		lease:    lease,
		JobName:  jobName,
		isLocked: false,
	}
}

// tryLock 尝试上锁
func (jobLock *JobLock) tryLock() error {
	// 1. 创建租约
	leaseGrantResp, err := jobLock.lease.Grant(context.TODO(), 5)
	if err != nil {
		return err
	}

	// 2. 自动续租
	leaseId := leaseGrantResp.ID
	ctx, cancelFunc := context.WithCancel(context.TODO())

	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc

	keepRespChan, err := jobLock.lease.KeepAlive(ctx, leaseId)
	if err != nil {
		if err := jobLock.revokeLease(); err != nil {
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

	// 3. 创建 TXN 事务
	// 获取事务
	txn := jobLock.kv.Txn(context.TODO())

	// 锁路径
	lockKey := common.JobLockDir + jobLock.JobName

	// 4. 事务抢锁
	// 定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		if err := jobLock.revokeLease(); err != nil {
			return err
		}
		return err
	}

	// 5. 成功返回, 失败释放租约
	// 抢锁失败
	if !txnResp.Succeeded {
		if err := jobLock.revokeLease(); err != nil {
			return err
		}
		return common.ErrorLockAlreadyRequired
	}

	// 抢锁成功
	jobLock.isLocked = true

	return nil
}

// 关闭续租并释放租约
func (jobLock *JobLock) revokeLease() error {
	// 关闭自动续租的协程
	jobLock.cancelFunc()

	// 释放租约
	if _, err := jobLock.lease.Revoke(context.TODO(), jobLock.leaseId); err != nil {
		return err
	}
	return nil
}

// unLock 释放锁
func (jobLock *JobLock) unLock() error {
	if jobLock.isLocked == true {
		if err := jobLock.revokeLease(); err != nil {
			return err
		}
	}
	return nil
}
