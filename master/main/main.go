package main

import (
	"cronTab/master/api_server"
	"cronTab/master/config"
	"cronTab/master/etcdOps"
	"flag"
	"github.com/golang/glog"
	"runtime"
	"time"
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
func initEnv() {
	// 设置线程数等于核心数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// 初始化命令行参数
	initArgs()

	defer glog.Flush()

	// 初始化线程
	initEnv()

	// 加载配置
	if err := config.InitialConfig(confFile); err != nil {
		glog.Fatal(err)
	}

	// 任务管理
	if err := etcdOps.InitJobMgr(); err != nil {
		glog.Fatal(err)
	}
	defer etcdOps.CloseEtcdConn()

	// 启动 Http 服务
	api_server.InitApiServer()

	for {
		time.Sleep(time.Second)
	}
}
