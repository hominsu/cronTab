package web_serve

import (
	"context"
	"cronTab/internal/master/biz/job_mgr"
	"cronTab/internal/master/service/log_sink"
	cron_job2 "cronTab/internal/pkg/cron_job"
	common2 "cronTab/internal/pkg/net"
	"cronTab/internal/pkg/xerrors"
	"github.com/gin-gonic/gin"
)

// 保存任务接口 POST job = {"name": "job1", "command": "echo hello", "cron_expr": "*/5 * * * * * *"}
func handlerJobSave(c *gin.Context) {
	// 任务保存到 etcd
	var err error

	// 反序列化 job
	job := &cron_job2.Job{}
	if err = c.BindJSON(job); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 保存到 etcd
	oldJob, err := job_mgr.GJobMgr.SaveJob(context.Background(), job)
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", oldJob)
}

// 保存任务接口 DELETE job = {"name": "job1", "command": "echo hello", "cron_expr": "*/5 * * * * * *"}
func handlerJobDelete(c *gin.Context) {
	// 删除 etcd 中的任务
	var err error

	// 反序列化 job
	job := &cron_job2.Job{}
	if err = c.BindJSON(job); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 删除任务
	oldJob, err := job_mgr.GJobMgr.DeleteJob(context.Background(), job.Name)
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", oldJob)
}

// 列举所有 crontab 任务
func handlerJobList(c *gin.Context) {
	var err error

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 获取任务
	jobList, err := job_mgr.GJobMgr.ListJobs(context.Background())
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", jobList)
}

// 强制杀死某个任务
func handlerJobKill(c *gin.Context) {
	var err error

	// 反序列化 job
	job := &cron_job2.Job{}
	if err = c.BindJSON(job); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 删除任务
	if err = job_mgr.GJobMgr.KillJob(context.Background(), job.Name); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", nil)
}

// 列出全部节点
func handlerNodeList(c *gin.Context) {
	var err error

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 列出节点
	nodes, err := job_mgr.GJobMgr.ListNodes(context.Background())
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", nodes)
}

// 获取任务日志
func handlerJobLogList(c *gin.Context) {
	var err error

	// 反序列化 JobPaging
	jobPaging := &cron_job2.JobPaging{}
	if err = c.BindJSON(jobPaging); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 获取任务日志
	logs, err := log_sink.GLogSink.GetLogBatch(context.Background(), jobPaging.Name)
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", logs)
}

// 删除任务日志
func handlerJobLogDelete(c *gin.Context) {
	var err error

	// 反序列化 JobPaging
	jobPaging := &cron_job2.JobPaging{}
	if err = c.BindJSON(jobPaging); err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	// 删除任务日志
	delCount, err := log_sink.GLogSink.DelJobLog(context.Background(), jobPaging.Name)
	if err != nil {
		common2.ResponseJson(c, -1, err.Error(), nil)
		xerrors.ErrFmt(err)
		return
	}

	// 返回正常应答
	common2.ResponseJson(c, 0, "success", delCount)
}
