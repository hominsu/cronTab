package main

import (
	"context"
	"cronTab/common"
	"cronTab/master/api_server"
	"cronTab/master/config"
	"cronTab/master/etcd_ops"
	"cronTab/master/job_mgr"
	"cronTab/master/mongodb_ops"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	flag.StringVar(&confFile, "config", "./master.json", "指定 master.json")
	flag.Parse()
}

func ginServe(handler *gin.Engine, stop chan struct{}) error {
	s := http.Server{
		Addr:    ":" + config.GConfig.ApiPort,
		Handler: handler,
	}

	go func() {
		<-stop
		if err := s.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	return s.ListenAndServe() // shutdown 之后会 return
}

func main() {
	var err error

	// 初始化命令行参数
	initArgs()

	// 初始化线程
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 加载配置
	if err = config.InitialConfig(confFile); err != nil {
		common.ErrFmtWithExit(err, 1)
	}

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

	errors := make(chan error, 1)
	stop := make(chan struct{})

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

	// 设置 gin 为 release 模式
	gin.SetMode(gin.ReleaseMode)

	// web 服务
	go func() {
		errors <- ginServe(api_server.HandlerRegister(), stop)
	}()

	// 阻塞等待退出
	<-done

	// 平滑退出
	close(stop)
	for i := 0; i < cap(errors); i++ {
		if err := <-errors; err != nil {
			log.Println(err)
		}
	}
}
