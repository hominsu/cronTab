package api_server

import (
	"cronTab/master/config"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type Engine struct {
	*gin.Engine
}

var (
	// GapiServer 单例对象
	GapiServer *Engine
)

func handlerRegister() (engine *gin.Engine) {
	e := Engine{gin.Default()}

	e.Handle("POST", "/job-save", handlerJobSave)

	return e.Engine
}

func InitApiServer() {
	// 设置 gin 为 release 模式
	//gin.SetMode(gin.ReleaseMode)

	defer glog.Flush()

	// 注册 api_server
	engine := handlerRegister()

	// 监听
	go func(engine *gin.Engine) {
		if err := engine.Run(":" + config.GConfig.ApiPort); err != nil {
			glog.Fatal(err)
		}
	}(engine)

	// 赋值单例模式
	GapiServer = &Engine{engine}

}
