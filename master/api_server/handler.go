package api_server

import (
	"cronTab/common"
	"cronTab/common/cron_job"
	"cronTab/master/job_mgr"
	"github.com/gin-gonic/gin"
)

// 保存任务接口 POST job = {"name": "job1", "command": "echo hello", "cron_expr": "*/5 * * * * * *"}
func handlerJobSave(c *gin.Context) {
	// 任务保存到 etcd
	var err error

	job := &cron_job.Job{}

	// 反序列化 job
	if err = c.BindJSON(job); err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 保存到 etcd
	oldJob, err := job_mgr.GJobMgr.SaveJob(job)
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
	var err error

	job := &cron_job.Job{}

	// 反序列化 job
	if err = c.BindJSON(job); err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 删除任务
	oldJob, err := job_mgr.GJobMgr.DeleteJob(job.Name)
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", oldJob)
}

// 列举所有 crontab 任务
func handlerJobList(c *gin.Context) {
	var err error

	// 获取任务
	jobList, err := job_mgr.GJobMgr.ListJobs()
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", jobList)
}

// 强制杀死某个任务
func handlerJobKill(c *gin.Context) {
	var err error

	job := &cron_job.Job{}

	// 反序列化 job
	if err = c.BindJSON(job); err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 删除任务
	if err = job_mgr.GJobMgr.KillJob(job.Name); err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", nil)
}

// 列出全部节点
func handlerNodeList(c *gin.Context) {
	// 列出节点
	nodes, err := job_mgr.GJobMgr.ListNodes()
	if err != nil {
		common.ResponseJson(c, -1, err.Error(), nil)
		return
	}

	// 返回正常应答
	common.ResponseJson(c, 0, "success", nodes)
}
