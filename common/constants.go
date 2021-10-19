package common

const (
	// JobSaveDir 任务保存目录
	JobSaveDir = "/cron/jobs/"

	// JobKillDir 任务强杀任务
	JobKillDir = "/cron/kill/"

	// JobEventSave 保存任务事件
	JobEventSave int = 1

	// JobEventDelete 删除任务事件
	JobEventDelete int = 2

	// JobLockDir 锁目录
	JobLockDir = "/cron/lock/"
)
