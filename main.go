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
	r.POST("/chat-process", gpt3client)

	if len(config.api_proxy) > 0 {
		r.Any("/api/*proxyPath", apiProxy)
	}
	// 设置public文件夹的静态文件服务
	r.Static("/assets", config.static_folder+"/assets")

	// //设置模板文件夹路径
	r.LoadHTMLGlob(config.static_folder + "/*.html")
	//设置路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run(fmt.Sprintf("%s:%s", config.http_host, config.http_port))
}
