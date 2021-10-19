package scheduler

import "cronTab/common"

// Scheduler 任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent               // etcd 任务队列
	jobPlanTable map[string]*common.JobSchedulerPlan // 任务调度计划表
}

var (
	GScheduler *Scheduler
)
