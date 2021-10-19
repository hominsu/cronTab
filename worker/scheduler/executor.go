package scheduler

import (
	"context"
	"cronTab/common"
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

		// 获取分布式锁

		// 任务结果
		result := &common.JobExecResult{
			ExecInfo:  info,
			Output:    make([]byte, 0), // 初始化一个空值
			StartTime: time.Now(),      // 记录任务开始时间
		}

		// 执行 shell 命令
		cmd := exec.CommandContext(ctx, "/bin/bash", "-c", info.Job.Command)

		// 执行并捕获输出
		output, err := cmd.CombinedOutput()

		// 记录任务结束时间、输出、错误
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err

		// 任务完成后, 把执行结果返回给 Scheduler, Scheduler 会从 executingTable 中删除执行记录
		GScheduler.PushJobResult(result)
	}()
}
