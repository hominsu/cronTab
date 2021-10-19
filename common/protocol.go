package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorhill/cronexpr"
	"net/http"
	"strings"
	"time"
)

// Job 定时任务
type Job struct {
	Name     string `json:"name"`      // 任务名
	Command  string `json:"command"`   // shell 命令
	CronExpr string `json:"cron_expr"` // cron 表达式
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
	return strings.TrimPrefix(jobKey, JobSaveDir)
}

// Response http 接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// ResponseJson 返回自定义格式的 Json
func ResponseJson(c *gin.Context, errno int, msg string, data interface{}) {
	// 返回应答
	//bytes, err := buildResponse(errno, msg, data)
	//if err != nil {
	//	c.String(http.StatusInternalServerError, "Internal Server Error")
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"errno": errno,
		"msg":   msg,
		"data":  data,
	})
}

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

// JobSchedulerPlan 任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 // 调度的任务信息
	Expr     *cronexpr.Expression // 解析好的 cron 表达式
	NextTime time.Time            // 下次调度时间
}

// BuildJobSchedulerPlan 构造执行计划
func BuildJobSchedulerPlan(job *Job) (*JobSchedulerPlan, error) {
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		return nil, err
	}

	// 生成任务调度计划对象
	return &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}, nil
}
