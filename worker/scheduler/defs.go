package scheduler

import "cronTab/common"

// Scheduler 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent               // etcd 任务队列
	jobPlanTable      map[string]*common.JobSchedulerPlan // 任务调度计划表
	jobExecutingTable map[string]*common.JobExecInfo      // 任务执行信息表
	jobResultChan     chan *common.JobExecResult          // 任务执行结果队列
}

var (
	GScheduler *Scheduler
)

// Executor 任务执行
type Executor struct {
}

var (
	GExecutor *Executor
)
