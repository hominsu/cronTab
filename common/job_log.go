package common

// JobLog 任务执行日志
type JobLog struct {
	JobName      string `bson:"job_name"`      // 任务名字
	Command      string `bson:"command"`       // 脚本命令
	Err          string `bson:"err"`           // 错误原因
	Output       string `bson:"output"`        // 脚本输出
	PlanTime     int64  `bson:"plan_time"`     // 计划开始时间
	ScheduleTime int64  `bson:"schedule_time"` // 实际调度时间
	StartTime    int64  `bson:"start_time"`    // 任务开始执行时间
	EndTime      int64  `bson:"end_time"`      // 任务结束执行时间
}

// LogBatch 日志批次
type LogBatch struct {
	Logs []interface{} // 多条日志
}
