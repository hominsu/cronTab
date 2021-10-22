package common

const (
	// JobSaveDir 任务保存目录
	JobSaveDir = "/cron/jobs/"

	// JobKillDir 任务强杀任务
	JobKillDir = "/cron/kill/"

	// JobLockDir 锁目录
	JobLockDir = "/cron/lock/"

	// JobEventSave 保存任务事件
	JobEventSave int = 1

	// JobEventDelete 删除任务事件
	JobEventDelete int = 2

	// JobEventKill 删除任务事件
	JobEventKill int = 3

	// NodeIpNet 节点 IP 地址
	NodeIpNet = "/cron/nodes/"
)
