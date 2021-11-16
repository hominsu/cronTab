package web_serve

import (
	"cronTab/configs/master_conf"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func HandlerRegister() *gin.Engine {
	engine := gin.Default()

	engine.Handle("POST", "/job", handlerJobSave)
	engine.Handle("DELETE", "/job", handlerJobDelete)
	engine.Handle("GET", "/job", handlerJobList)
	engine.Handle("POST", "/job/kill", handlerJobKill)
	engine.Handle("GET", "/job/node", handlerNodeList)
	engine.Handle("POST", "/job/log", handlerJobLogList)
	engine.Handle("DELETE", "/job/log", handlerJobLogDelete)

	// 静态
	engine.Static("/static", master_conf.GConfig.WebRoot)

	engine.LoadHTMLGlob(path.Join(master_conf.GConfig.WebRoot, "html/*"))
	engine.Handle("GET", "/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return engine
}
