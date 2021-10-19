package main

import (
	"cronTab/worker/config"
	"cronTab/worker/etcdOps"
	"cronTab/worker/scheduler"
	"flag"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "指定 worker.json")
	flag.Parse()
}

// 初始化线程
func initEnv() {
	// 设置线程数等于核心数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var err error

	// 初始化命令行参数
	initArgs()

	defer glog.Flush()

	done := make(chan bool)
	// 创建监听 chan
	ch := make(chan os.Signal)
	// 监听退出信号
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 开一个 goroutine 捕获退出信号
	go func(done chan bool) {
		for s := range ch {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				glog.Info("Program Exit...", s)
				done <- true
			default:
				glog.Warning("Other signal", s)
			}
		}
	}(done)

	// 初始化线程
	initEnv()

	// 加载配置
	if err = config.InitialConfig(confFile); err != nil {
		glog.Fatal(err)
	}

	// 任务管理
	if err = etcdOps.InitJobMgr(); err != nil {
		glog.Fatal(err)
	}
	defer etcdOps.CloseEtcdConn()

	// 启动执行器
	if err = scheduler.InitExecutor(); err != nil {
		glog.Fatal(err)
	}

	// 启动调度
	if err = scheduler.InitScheduler(); err != nil {
		glog.Fatal(err)
	}

	// 启动监听
	if err = etcdOps.GJobMgr.WatchJob(); err != nil {
		glog.Fatal(err)
	}

	// 阻塞，等待退出
	<-done
}
