package main

import (
	"cronTab/common"
	"cronTab/master/api_server"
	"cronTab/master/config"
	"cronTab/master/etcd_ops"
	"cronTab/master/job_mgr"
	"cronTab/master/mongodb_ops"
	"flag"
	"github.com/golang/glog"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "指定 master.json")
	flag.Parse()
}

// 初始化线程
func initProcess() {
	// 设置线程数等于核心数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var err error

	// 初始化命令行参数
	initArgs()

	defer glog.Flush()

	// 初始化线程
	initProcess()

	// 加载配置
	if err = config.InitialConfig(confFile); err != nil {
		common.ErrFmtWithExit(err, 1)
	}

	// 连接 etcd
	if err = etcd_ops.InitEtcdConn(); err != nil {
		common.ErrFmtWithExit(err, 1)
	}
	defer func() {
		err := etcd_ops.CloseEtcdConn()
		if err != nil {
			common.ErrFmt(err)
		}
	}()

	// 连接 mongodb
	if err = mongodb_ops.InitMongodbConn(); err != nil {
		common.ErrFmtWithExit(err, 1)
	}
	defer func() {
		err := mongodb_ops.CloseMongodbConn()
		if err != nil {
			common.ErrFmt(err)
		}
	}()

	// 任务管理
	if err = job_mgr.InitJobMgr(); err != nil {
		common.ErrFmtWithExit(err, 1)
	}

	// 启动 Http 服务
	if err = api_server.InitApiServer(); err != nil {
		common.ErrFmtWithExit(err, 1)
	}
}
