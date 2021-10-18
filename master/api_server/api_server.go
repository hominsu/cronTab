package api_server

import (
	"cronTab/master/config"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net/http"
	"path"
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

	e.Handle("POST", "/job/save", handlerJobSave)
	e.Handle("POST", "/job/delete", handlerJobDelete)
	e.Handle("GET", "/job/list", handlerJobList)
	e.Handle("POST", "/job/kill", handlerJobKill)

	// 静态
	e.Static("/static", config.GConfig.WebRoot)
	//fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./web")))
	//handler := func(c *gin.Context) {
	//	fileServer.ServeHTTP(c.Writer, c.Request)
	//}
	//e.GET("/*filepath", handler)

	e.LoadHTMLGlob(path.Join(config.GConfig.WebRoot, "html/*"))
	e.Handle("GET", "/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

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
