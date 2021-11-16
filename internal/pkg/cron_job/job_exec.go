package cron_job

import (
	"context"
	"time"
)

// JobExecInfo 任务执行状态
type JobExecInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 计划调度时间
	RealTime   time.Time          // 实际调度时间
	CancelCtx  context.Context    // 任务 command 的上下文
	CancelFunc context.CancelFunc // 用于取消 command 执行的函数
}

// BuildJobExecInfo 构造执行状态信息
func BuildJobExecInfo(jobSchedulerPlan *JobSchedulerPlan) *JobExecInfo {
	jobExecInfo := &JobExecInfo{
		Job:      jobSchedulerPlan.Job,
		PlanTime: jobSchedulerPlan.NextTime, // 计划调度时间
		RealTime: time.Now(),                //真实调度时间
	}

	jobExecInfo.CancelCtx, jobExecInfo.CancelFunc = context.WithCancel(context.Background())

	return jobExecInfo
}

// JobExecResult 任务执行结果
type JobExecResult struct {
	ExecInfo  *JobExecInfo // 执行状态
	Output    []byte       // 脚本数据
	Err       error        // 脚本错误原因
	StartTime time.Time    // 开始时间
	EndTime   time.Time    // 结束时间
}
