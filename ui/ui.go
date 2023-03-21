package ui

import (
	"net/http"
	"os"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

var WebFileServer http.Handler = http.FileServer(ChatGptWeb())

func ChatGptWeb() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		k = "web"
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k, Fallback: "index.html"}
	}
	panic("unreachable")
}

func BinStatic(c *gin.Context) {

	length := len(c.Request.URL.Path)

	if (length >= 5 && c.Request.URL.Path[0:5] == "/api/") ||
		(length >= 11 && c.Request.URL.Path[0:11] == "/websocket/") { // API & websocket
		println("api called,return")
		c.Next()
		return
	}
	WebFileServer.ServeHTTP(c.Writer, c.Request)
	c.Abort()

}
