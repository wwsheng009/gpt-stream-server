package main

import (
	"fmt"
	"gpt_stream_server/chat"
	"gpt_stream_server/config"
	"gpt_stream_server/gpt"
	"gpt_stream_server/ui"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	r := gin.Default()
	r.Use(ui.BinStatic)

	// Simple group: v2
	api := r.Group("/api")
	{
		// api.POST("/stream", gpt.StreamHandler)
		api.POST("/chat-process", chat.Auth, gpt.StreamHandler)
		api.POST("/session", chat.Session)
		api.POST("/config", chat.Auth, chat.ConfigHandler)
		api.POST("/verify", chat.Verify)
		api.Any("/yao/*proxyPath", gpt.ApiProxy)
	}

	// if len(config.MainConfig.ApiServer) > 0 {
	//
	// }

	// 设置public文件夹的静态文件服务
	//r.Static("/assets", config.MainConfig.StaticFolder+"/assets")

	// //设置模板文件夹路径
	//r.LoadHTMLGlob(config.MainConfig.StaticFolder + "/*.html")
	//设置路由

	// r.Use(static.Serve("/assets", ui.BinaryFileSystem("web/assets")))
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{})
	// })

	r.Run(fmt.Sprintf("%s:%s", config.MainConfig.HttpHost, config.MainConfig.HttpPort))
}
