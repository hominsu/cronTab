package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response http 接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// ResponseJson 返回自定义格式的 Json
func ResponseJson(c *gin.Context, errno int, msg string, data interface{}) {
	// 返回应答
	//bytes, err := buildResponse(errno, msg, data)
	//if err != nil {
	//	c.String(http.StatusInternalServerError, "Internal Server Error")
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"errno": errno,
		"msg":   msg,
		"data":  data,
	})
}
