package scheduler

import (
	"cronTab/common"
	"time"
)

// InitScheduler 初始化调度器
func InitScheduler() error {
	GScheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*common.JobExecInfo),
		jobResultChan:     make(chan *common.JobExecResult, 1000),
	}

	// 启动调度协程
	go GScheduler.schedulerLoop()

	return nil
}

// PushJobEvent 推送任务事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// 调度协程
func (scheduler *Scheduler) schedulerLoop() {
	// 初始化调度的延时定时器
	schedulerTimer := time.NewTimer(scheduler.TryScheduler())

	for {
		select {
		case jobEvent := <-GScheduler.jobEventChan: // 监听任务变化事件
			// 对内存中维护的任务列表做增删改查
			scheduler.handlerJobEvent(jobEvent)
		case <-schedulerTimer.C: // 任务到期
		case jobResult := <-scheduler.jobResultChan: // 监听任务执行结果
			scheduler.handlerJobResult(jobResult)
		}
		// 调度一次任务并重置调度间隔
		schedulerTimer.Reset(scheduler.TryScheduler())
	}
}

// 处理任务更改
func (scheduler *Scheduler) handlerJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	// 保存事件
	case common.JobEventSave:
		jobSchedulerPlan, err := common.BuildJobSchedulerPlan(jobEvent.Job)
		if err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	// 删除事件
	case common.JobEventDelete:
		if _, ok := scheduler.jobPlanTable[jobEvent.Job.Name]; ok {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// TryScheduler 尝试调度任务并重新计算任务调度状态
func (scheduler Scheduler) TryScheduler() time.Duration {
	var nearTime *time.Time

	// 如果任务表为空, 随便睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		return time.Second
	}

	// 当前时间
	now := time.Now()

	// 1. 遍历所有任务
	for _, jobPlan := range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			// 尝试执行任务
			scheduler.TryStartJob(jobPlan)

			// 更新下次执行时间
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}

		// 统计最近一个要过期的时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	return (*nearTime).Sub(now)
}

// TryStartJob 尝试执行任务
func (scheduler Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	// 如果调度的时间间隔小于任务执行所需时间, 只能执行一次, 防止并发

	// 如果任务正在执行, 跳过本次调度
	if _, ok := scheduler.jobExecutingTable[jobPlan.Job.Name]; ok {
		//fmt.Println("跳过:", jobPlan.Job.Name)
		return
	}

	// 构建执行状态信息
	jobExecInfo := common.BuildJobExecInfo(jobPlan)

	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecInfo

	// 执行任务
	GExecutor.ExecJob(jobExecInfo)

	//bytes, _ := jobPlan.Job.JobMarshal()
	//fmt.Printf("%s, [%s], [%s]\n", string(bytes), jobExecInfo.PlanTime, jobExecInfo.RealTime)
}

// PushJobResult 回传任务执行结果
func (scheduler Scheduler) PushJobResult(jobResult *common.JobExecResult) {
	scheduler.jobResultChan <- jobResult
}

// 处理任务结果
func (scheduler Scheduler) handlerJobResult(result *common.JobExecResult) {
	// 删除任务执行状态
	if _, ok := scheduler.jobExecutingTable[result.ExecInfo.Job.Name]; ok {
		delete(scheduler.jobExecutingTable, result.ExecInfo.Job.Name)
		//fmt.Println(result.ExecInfo.Job.Name, strings.Replace(string(result.Output), "\n", "", -1), result.Err)
	}
}
