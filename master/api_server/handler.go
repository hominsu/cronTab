package api_server

import (
	"cronTab/common"
	"cronTab/master/etcdOps"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// 保存任务接口 POST job = {"name": "job1", "command": "echo hello", "cron_expr": "*/5 * * * * * *"}
func handlerJobSave(c *gin.Context) {
	// 任务保存到 etcd

	// 获取表单中 job 的值
	postJob := c.PostForm("job")

	// 反序列化 job
	job := &common.Job{}
	if err := json.Unmarshal([]byte(postJob), job); err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 保存到 etcd
	oldJob, err := etcdOps.GJobMgr.SaveJob(job)
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", oldJob)
}

// 保存任务接口 DELETE job = {"name": "job1", "command": "echo hello", "cron_expr": "*/5 * * * * * *"}
func handlerJobDelete(c *gin.Context) {
	// 删除 etcd 中的任务

	// 获取表单中 name 的值
	name := c.PostForm("name")

	// 删除任务
	oldJob, err := etcdOps.GJobMgr.DeleteJob(name)
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", oldJob)
}

// 列举所有 crontab 任务
func handlerJobList(c *gin.Context) {

	// 获取任务
	jobList, err := etcdOps.GJobMgr.ListJobs()
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", jobList)
}

// 强制杀死某个任务
func handlerJobKill(c *gin.Context) {
	// 获取表单中 name 的值
	name := c.PostForm("name")

	// 删除任务
	err := etcdOps.GJobMgr.KillJob(name)
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", nil)
}
