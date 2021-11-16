package cron_job

// JobLog 任务执行日志
type JobLog struct {
	JobName      string `bson:"job_name" json:"job_name"`           // 任务名字
	Command      string `bson:"command" json:"command"`             // 脚本命令
	Err          string `bson:"err" json:"err"`                     // 错误原因
	Output       string `bson:"output" json:"output"`               // 脚本输出
	PlanTime     int64  `bson:"plan_time" json:"plan_time"`         // 计划开始时间
	ScheduleTime int64  `bson:"schedule_time" json:"schedule_time"` // 实际调度时间
	StartTime    int64  `bson:"start_time" json:"start_time"`       // 任务开始执行时间
	EndTime      int64  `bson:"end_time" json:"end_time"`           // 任务结束执行时间
}

// JobPaging 任务分页获取
type JobPaging struct {
	Name  string `json:"name,omitempty"`  // 任务名字
	Skip  int64  `json:"skip,omitempty"`  // 跳过多少条
	Limit int64  `json:"limit,omitempty"` // 限制多少条
}

// LogBatch 日志批次
type LogBatch struct {
	Logs []interface{} // 多条日志
}
