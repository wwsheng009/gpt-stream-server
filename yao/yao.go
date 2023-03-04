package yao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gpt_stream_server/config"
	"io"
	"net/http"
)

type PayloadRequest struct {
	Type    string      `json:"type"`
	Method  string      `json:"method"`
	Args    interface{} `json:"args,omitempty"`
	Space   string      `json:"space,omitempty"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	Message string      `json:"message,omitempty"`
}

func RemoteRequest(payload PayloadRequest) (interface{}, error) {
	if config.MainConfig.ApiServer == "" {
		return nil, errors.New("代理地址为空，请配置环境变量YAO_APP_PROXY_ENDPOINT")
	}
	url := config.MainConfig.ApiServer + "/proxy/call"
	// println("url:", url)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.MainConfig.ApiServerAccessKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Response status code is %d\n", resp.StatusCode)
		// Handle any errors or return accordingly
		res, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("远程程序执行异常:代码:%d,消息：:%s", resp.StatusCode, string(res))
	}

	defer resp.Body.Close()
	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	if res["message"] != "" && res["code"] != 200 {
		return nil, fmt.Errorf("远程程序执行异常:代码:%d,消息：%s", resp.StatusCode, res["message"])
	}
	return res["data"], nil
}

func YaoProcess(method string, args ...interface{}) (interface{}, error) {
	return RemoteRequest(PayloadRequest{
		Type:   "Process",
		Method: method,
		Args:   args,
	})
}
