package main

import (
	"fmt"
	"gpt_stream_server/config"
	"gpt_stream_server/gpt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	r := gin.Default()

	r.POST("/stream", gpt.StreamHandler)
	r.POST("/chat-process", gpt.StreamHandler)

	if len(config.MainConfig.ApiServer) > 0 {
		r.Any("/api/*proxyPath", gpt.ApiProxy)
	}
	// 设置public文件夹的静态文件服务
	r.Static("/assets", config.MainConfig.StaticFolder+"/assets")

	// //设置模板文件夹路径
	r.LoadHTMLGlob(config.MainConfig.StaticFolder + "/*.html")
	//设置路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run(fmt.Sprintf("%s:%s", config.MainConfig.HttpHost, config.MainConfig.HttpPort))
}
