package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Response http 接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 构建应答 Json
func buildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 定义 Response 对象
	response := &Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}

	// 序列化
	resp, err = json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
