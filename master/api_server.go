package master

import (
	"cronTab/master/config"
	"flag"
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

	e.Handle("DELETE", "/logout", handlerJobSave)

	return e.Engine
}

func InitApiServer() {
	// 设置 gin 为 release 模式
	//gin.SetMode(gin.ReleaseMode)

	flag.Parse()
	defer glog.Flush()

	// 注册 handler
	engine := handlerRegister()

	// 监听
	go func(engine *gin.Engine) {
		if err := engine.Run(":" + config.Basic.WebPort()); err != nil {
			glog.Fatal(err)
		}
	}(engine)

	// 赋值单例模式
	GapiServer = &Engine{engine}

}
