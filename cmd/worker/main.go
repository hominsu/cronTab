package main

import (
	"context"
	"cronTab/configs/worker_conf"
	"cronTab/internal/pkg/xerrors"
	"cronTab/internal/worker/biz/job_mgr"
	"cronTab/internal/worker/data/etcd"
	"cronTab/internal/worker/data/mongodb"
	"cronTab/internal/worker/service/heart_beat"
	"cronTab/internal/worker/service/log_sink"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "指定 worker.json")
	flag.Parse()
}

func main() {
	var err error

	// 初始化命令行参数
	initArgs()

	done := make(chan bool)
	ch := make(chan os.Signal)

	// 监听退出信号
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(done chan bool) {
		for s := range ch {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Program Exit...", s)
				done <- true
			default:
				log.Println("Other signal", s)
			}
		}
	}(done)

	errors := make(chan error, 4)
	stop := make(chan struct{})

	// 初始化线程
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置
	if err = worker_conf.InitialConfig(confFile); err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 连接 etcd
	if err = etcd.InitEtcdConn(ctx); err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}
	defer func() {
		err := etcd.CloseEtcdConn()
		if err != nil {
			xerrors.ErrFmt(err)
		}
	}()

	// 连接 mongodb
	if err = mongodb.InitMongodbConn(ctx); err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}
	defer func() {
		err := mongodb.CloseMongodbConn()
		if err != nil {
			xerrors.ErrFmt(err)
		}
	}()

	// 初始化心跳
	if err = heart_beat.InitHeartBeat(ctx); err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}
	defer func() {
		if err := heart_beat.StopHeartBeat(); err != nil {
			xerrors.ErrFmt(err)
		}
	}()

	// 启动日志池
	if err = log_sink.InitLogSink(); err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}

	// 启动日志池处理协程
	go func() {
		errors <- log_sink.GLogSink.WriteLoop(stop)
	}()

	// 初始化任务管理
	jobMgr, err := job_mgr.InitJobMgr()
	if err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}

	// 启动任务调度协程
	go func() {
		errors <- job_mgr.GScheduler.SchedulerLoop(stop)
	}()

	// 启动任务监听协程
	jobRevision, err := jobMgr.GetJobRevision(ctx)
	if err != nil {
		xerrors.ErrFmtWithExit(err, 1)
	}
	go func() {
		errors <- jobMgr.WatchJob(jobRevision, stop)
	}()

	// 启动强杀监听协程
	go func() {
		errors <- jobMgr.WatchKill(stop)
	}()

	// 阻塞，等待退出
	<-done

	// 平滑退出
	close(stop)
	for i := 0; i < cap(errors); i++ {
		if err := <-errors; err != nil {
			xerrors.ErrFmt(err)
		}
	}
}
