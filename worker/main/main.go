package main

import (
	"cronTab/worker/config"
	"cronTab/worker/etcd_ops"
	"cronTab/worker/heart_beat"
	"cronTab/worker/job_mgr"
	"cronTab/worker/mongodb_ops"
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

	// 连接 etcd
	if err = etcd_ops.InitEtcdConn(); err != nil {
		glog.Fatal(err)
	}
	defer func() {
		err := etcd_ops.CloseEtcdConn()
		if err != nil {
			glog.Fatal(err)
		}
	}()

	// 连接 mongodb
	if err = mongodb_ops.InitMongodbConn(); err != nil {
		glog.Fatal(err)
	}
	defer func() {
		err := mongodb_ops.CloseMongodbConn()
		if err != nil {
			glog.Fatal(err)
		}
	}()

	// 初始化心跳
	heartBeat := heart_beat.InitHeartBeat()
	if err = heartBeat.StartHeartBeat(); err != nil {
		glog.Fatal(err)
	}
	defer func(heartBeat *heart_beat.HeartBeat) {
		err := heartBeat.EndHeartBeat()
		if err != nil {
			glog.Fatal(err)
		}
	}(heartBeat)

	// 初始化任务管理
	if err = job_mgr.InitJobMgr(); err != nil {
		glog.Fatal(err)
	}

	// 阻塞，等待退出
	<-done
}
