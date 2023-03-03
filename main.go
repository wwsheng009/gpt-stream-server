package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	Init_config()

	r := gin.Default()

	r.POST("/stream", streamHandler)

	if len(config.api_proxy) > 0 {
		r.Any("/api/*proxyPath", apiProxy)
	}
	// 设置public文件夹的静态文件服务
	r.Static("/public", "./public")

	// //设置模板文件夹路径
	// r.LoadHTMLGlob("views/*")
	//设置路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run(fmt.Sprintf("%s:%s", config.http_host, config.http_port))
}

func streamHandler(c *gin.Context) {

	// 设置响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	// 读取请求体

	if c.Request.Method == "POST" {
		// 解析JSON
		var jsonBody map[string]interface{}
		err := c.ShouldBindJSON(jsonBody)
		if err != nil {
			// 处理错误
			return
		}
		if _, has := jsonBody["message"]; !has {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "message is empty"})
			return
		}
	}

	processComplete(c.Writer, c.Request)
	// 处理请求
}
