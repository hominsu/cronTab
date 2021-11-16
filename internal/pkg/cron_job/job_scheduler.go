package cron_job

import (
	"github.com/gorhill/cronexpr"
	"time"
)

// JobSchedulerPlan 任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 // 调度的任务信息
	Expr     *cronexpr.Expression // 解析好的 cron 表达式
	NextTime time.Time            // 下次调度时间
}

// BuildJobSchedulerPlan 构造执行计划
func BuildJobSchedulerPlan(job *Job) (*JobSchedulerPlan, error) {
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		return nil, err
	}

	// 生成任务调度计划对象
	return &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}, nil
}
