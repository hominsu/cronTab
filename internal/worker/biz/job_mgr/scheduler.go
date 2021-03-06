package job_mgr

import (
	"cronTab/internal/pkg/constants"
	"cronTab/internal/pkg/cron_job"
	"cronTab/internal/pkg/xerrors"
	"cronTab/internal/worker/service/log_sink"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Scheduler 任务调度
type Scheduler struct {
	jobEventChan      chan *cron_job.JobEvent               // etcd 任务队列
	jobPlanTable      map[string]*cron_job.JobSchedulerPlan // 任务调度计划表
	jobExecutingTable map[string]*cron_job.JobExecInfo      // 任务执行信息表
	jobResultChan     chan *cron_job.JobExecResult          // 任务执行结果队列
}

var (
	GScheduler *Scheduler
)

// InitScheduler 初始化调度器
func InitScheduler() error {
	GScheduler = &Scheduler{
		jobEventChan:      make(chan *cron_job.JobEvent, 1000),
		jobPlanTable:      make(map[string]*cron_job.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*cron_job.JobExecInfo),
		jobResultChan:     make(chan *cron_job.JobExecResult, 1000),
	}

	return nil
}

// PushJobEvent 推送任务事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *cron_job.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

// SchedulerLoop 调度协程
func (scheduler *Scheduler) SchedulerLoop(stop <-chan struct{}) error {
	// 初始化调度的延时定时器
	schedulerTimer := time.NewTimer(scheduler.tryScheduler())

	for {
		select {
		case <-schedulerTimer.C: // 任务到期
		case jobEvent := <-GScheduler.jobEventChan: // 监听任务变化事件
			scheduler.handlerJobEvent(jobEvent)
		case jobResult := <-scheduler.jobResultChan: // 监听任务执行结果
			scheduler.handlerJobResult(jobResult)
		case <-stop:
			// 取消正在执行的任务
			for _, jobExecInfo := range scheduler.jobExecutingTable {
				jobExecInfo.CancelFunc()
				delete(scheduler.jobExecutingTable, jobExecInfo.Job.Name)
			}
			return nil
		}
		// 调度一次任务并重置调度间隔
		schedulerTimer.Reset(scheduler.tryScheduler())
	}
}

// 处理任务更改
func (scheduler *Scheduler) handlerJobEvent(jobEvent *cron_job.JobEvent) {
	switch jobEvent.EventType {
	case constants.JobEventSave: // 保存事件
		jobSchedulerPlan, err := cron_job.BuildJobSchedulerPlan(jobEvent.Job)
		if err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	case constants.JobEventDelete: // 删除事件
		if _, ok := scheduler.jobPlanTable[jobEvent.Job.Name]; ok {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case constants.JobEventKill: // 强杀事件
		if jobExecInfo, ok := scheduler.jobExecutingTable[jobEvent.Job.Name]; ok {
			// 触发 command 杀死 shell 子进程
			jobExecInfo.CancelFunc()
		}
	}
}

// 处理任务结果
func (scheduler *Scheduler) handlerJobResult(result *cron_job.JobExecResult) {
	// 删除任务执行状态
	if _, ok := scheduler.jobExecutingTable[result.ExecInfo.Job.Name]; ok {
		delete(scheduler.jobExecutingTable, result.ExecInfo.Job.Name)
		fmt.Printf("%s: output:[%s], err:[%s]\n", result.ExecInfo.Job.Name, strings.Replace(string(result.Output), "\n", "", -1), result.Err)
	}

	// 生成执行日志
	//if result.Err != common.ErrorLockAlreadyRequired {
	if !errors.Is(result.Err, xerrors.ErrorLockAlreadyRequired) {
		jobLog := &cron_job.JobLog{
			JobName:      result.ExecInfo.Job.Name,
			Command:      result.ExecInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}

		// 存储到 mongodb
		log_sink.GLogSink.Append(jobLog)
	}
}

// 尝试调度任务并重新计算任务调度状态
func (scheduler *Scheduler) tryScheduler() time.Duration {
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
			scheduler.tryStartJob(jobPlan)

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

// 尝试执行任务
func (scheduler *Scheduler) tryStartJob(jobPlan *cron_job.JobSchedulerPlan) {
	// 如果调度的时间间隔小于任务执行所需时间, 只能执行一次, 防止并发

	// 如果任务正在执行, 跳过本次调度
	if _, ok := scheduler.jobExecutingTable[jobPlan.Job.Name]; ok {
		//fmt.Println("跳过:", jobPlan.Job.Name)
		return
	}

	// 构建执行状态信息
	jobExecInfo := cron_job.BuildJobExecInfo(jobPlan)

	// 保存执行状态
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecInfo

	// 执行任务
	GExecutor.ExecJob(jobExecInfo)

	bytes, _ := jobPlan.Job.JobMarshal()
	fmt.Printf("%s, [%s], [%s]\n", string(bytes), jobExecInfo.PlanTime, jobExecInfo.RealTime)
}

// PushJobResult 回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *cron_job.JobExecResult) {
	scheduler.jobResultChan <- jobResult
}
