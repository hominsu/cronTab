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

func handlerRegister() (engine *gin.Engine) {
	e := Engine{gin.Default()}

	e.Handle("POST", "/job", handlerJobSave)
	e.Handle("DELETE", "/job", handlerJobDelete)
	e.Handle("GET", "/job", handlerJobList)
	e.Handle("POST", "/job/kill", handlerJobKill)
	e.Handle("GET", "/job/node", handlerNodeList)

	// 静态
	e.Static("/static", config.GConfig.WebRoot)

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
	if err := engine.Run(":" + config.GConfig.ApiPort); err != nil {
		glog.Fatal(err)
	}

}
