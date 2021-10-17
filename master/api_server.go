package master

import (
	"log"
	"net"
	"net/http"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	// G_apiServer 单例对象
	G_apiServer *ApiServer
)

func handlerJobSave(w http.ResponseWriter, r *http.Request) {

}

func initApiServer() (err error) {
	// 配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handlerJobSave)

	// 启动 TCP 监听
	listener, err := net.Listen("tcp", ":8070")
	if err != nil {
		return err
	}

	// 创建 http 服务
	httpServer := &http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动服务端
	go func() {
		err := httpServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}
