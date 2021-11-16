package cron_job

// JobEvent 任务变化事件: 更新任务, 删除任务
type JobEvent struct {
	EventType int // SAVE, DELETE
	Job       *Job
}

// BuildJobEvent 构建 Event
func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}
