package scheduler

import (
	"context"
	"cronTab/common"
	"cronTab/worker/etcdOps"
	"github.com/golang/glog"
	"os/exec"
	"time"
)

// InitExecutor 初始化执行器
func InitExecutor() error {
	GExecutor = &Executor{}

	return nil
}

// ExecJob 执行任务
func (executor Executor) ExecJob(info *common.JobExecInfo) {
	go func() {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		// 初始化分布式锁
		jobLock := etcdOps.GJobMgr.CreateJobLock(info.Job.Name)

		// 任务结果
		result := &common.JobExecResult{
			ExecInfo:  info,
			Output:    make([]byte, 0), // 初始化一个空值
			StartTime: time.Now(),      // 记录任务开始时间
		}

		// 尝试上锁
		err := jobLock.TryLock()
		// 释放锁
		defer func(jobLock *etcdOps.JobLock) {
			if err := jobLock.UnLock(); err != nil {
				glog.Warning(err)
			}
		}(jobLock)

		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// 上锁成功, 重置任务开始时间
			result.StartTime = time.Now()

			// 执行 shell 命令
			cmd := exec.CommandContext(ctx, "/bin/bash", "-c", info.Job.Command)

			// 执行并捕获输出
			output, err := cmd.CombinedOutput()

			// 记录任务结束时间、输出、错误
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}

		// 任务完成后, 把执行结果返回给 Scheduler, Scheduler 会从 executingTable 中删除执行记录
		GScheduler.PushJobResult(result)
	}()
}
