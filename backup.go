package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func _processChat_backup(w http.ResponseWriter, r *http.Request) {
	// 创建一个HTTP客户端
	client := &http.Client{}

	var content = "use js to write a program post user login info to background api server"
	// content = "讲出你的名字"
	request := ChatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []Message{{Role: "system", Content: "you are assistant"}, {Role: "user", Content: content}},
		Stream:   true,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}

	// 发送另一个Stream API的请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", &buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Accept", "text/event-stream; charset=utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.openai_key)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 从另一个Stream API的响应中读取数据并发送给客户端

	lastline := ""
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			buf := make([]byte, 1024)
			n, err := resp.Body.Read(buf)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			if n > 0 {
				buf := buf[:n]
				data := lastline + string(buf)
				lastline = ""
				if data == "[DONE]" {
					return
				}
				lines := strings.Split(data, "\n\n")
				if len(lines) == 0 {
					lines = append(lines, data)
				}
				for i := 0; i < len(lines); i++ {
					if lines[i] == "" {

						continue
					}
					lines[i] = strings.TrimPrefix(lines[i], "data:")
					lines[i] = strings.Trim(lines[i], " ")
					if lines[i][len(lines[i])-1] != '}' {
						lastline = lines[i]
						if lastline == "[DONE]" {
							return
						}
						break
					}
					var res ChatResponse
					err := json.Unmarshal([]byte(lines[i]), &res)
					if err != nil {
						lastline = lines[i]
						// fmt.Println("Error reading response:", err)
						break
					}
					if res.Choices[0].Delta.Content == "" && res.Choices[0].FinishReason == "stop" {
						return
					}
					fmt.Fprintf(w, "%s", res.Choices[0].Delta.Content)
					w.(http.Flusher).Flush()
				}
			}
		}
	}
}
