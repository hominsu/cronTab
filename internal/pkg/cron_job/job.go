package cron_job

import (
	"cronTab/internal/pkg/constants"
	"encoding/json"
	"strings"
)

// Job 定时任务
type Job struct {
	Name     string `json:"name,omitempty"`      // 任务名
	Command  string `json:"command,omitempty"`   // shell 命令
	CronExpr string `json:"cron_expr,omitempty"` // cron 表达式
}

// JobMarshal 序列化 Job
func (j Job) JobMarshal() ([]byte, error) {
	if bytes, err := json.Marshal(j); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

// JobUnmarshal 反序列化 Job
func JobUnmarshal(bytes []byte) (*Job, error) {
	job := &Job{}
	// 对job进行反序列化
	if err := json.Unmarshal(bytes, job); err != nil {
		return nil, err
	}
	return job, nil
}

// ExtractJobName 从 etcd 中的 key 提取任务名
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, constants.JobSaveDir)
}

// ExtractKillName 从 etcd 中的 key 提取任务名
func ExtractKillName(killKey string) string {
	return strings.TrimPrefix(killKey, constants.JobKillDir)
}
