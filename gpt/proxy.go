package gpt

import (
	"bytes"
	"gpt_stream_server/config"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// request another http resource in gin route
func _http_request(c *gin.Context) {
	// Create a new request to the target server
	targetURL := "https://example.com"
	req, err := http.NewRequest("GET", targetURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers and other options as needed
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body and send it back to the client
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// pass through the api request to api server
func ApiProxy(c *gin.Context) {
	path := c.Param("proxyPath")
	if path == "/chat-process" && c.Request.Method == "POST" {
		StreamHandler(c)
		return
	}
	remote, err := url.Parse(config.MainConfig.ApiServer)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	originalDirector := proxy.Director

	// 修改请求路径
	// c.Request.URL.Path = c.Param("proxyPath")

	//Define the director func

	//This is a good place to log, for example

	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header = c.Request.Header
	}

	proxy.ModifyResponse = modifyResponse
	proxy.ErrorHandler = errorHandler
	proxy.ServeHTTP(c.Writer, c.Request)

}

func errorHandler(http.ResponseWriter, *http.Request, error) {
	println("error handle")
}
func modifyResponse(*http.Response) error {
	return nil
}
