package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func processComplete(w http.ResponseWriter, r *http.Request) {
	var dataPrefix = []byte("data: ")
	var doneSequence = []byte("[DONE]")
	// 创建一个HTTP客户端
	client := &http.Client{}

	var tokens = 1024
	var content = "use js to write a program post user login info to background api server"
	// content = "讲出你的名字"
	request := CompletionRequest{
		Model:     "text-davinci-003",
		Prompt:    []string{content},
		MaxTokens: &tokens,
		Stream:    true,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		log.Fatal(err)
	}

	// 发送另一个Stream API的请求
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", &buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Accept", "text/event-stream; charset=utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+"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")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// 从另一个Stream API的响应中读取数据并发送给客户端

	// lastline := []byte{}
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			line, err := reader.ReadSlice('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				writeError(w, err)
				return
			}
			if len(line) == 1 && line[0] == '\n' {
				continue
			}
			// if line[len(line)-2] != '}' {
			// 	lastline = line
			// 	continue
			// }
			// if len(lastline) > 0 {
			// 	line = append(lastline, line...)
			// 	lastline = []byte{}
			// }
			// make sure there isn't any extra whitespace before or after
			line = bytes.TrimSpace(line)
			// the completion API only returns data events
			if !bytes.HasPrefix(line, dataPrefix) {
				continue
			}
			line = bytes.TrimPrefix(line, dataPrefix)

			// the stream is completed when terminated by [DONE]
			if bytes.HasPrefix(line, doneSequence) {
				writeDone(w)
				break
			}

			var res CompletionResponse
			err = json.Unmarshal(line, &res)
			if err != nil {
				// lastline = line
				// fmt.Println("Error reading response:", err)
				continue
			}
			if res.Choices[0].Text == "" && res.Choices[0].FinishReason == "stop" {
				writeDone(w)
				return
			}
			writeMessage(w, res.Choices[0].Text)
			// fmt.Fprintf(w, `{"code":200,"message":"%s"}`, res.Choices[0].Text)
			// w.(http.Flusher).Flush()
			// c.JSON(200, gin.H{"message": res.Choices[0].Text})
		}
	}

}

func processChat(w http.ResponseWriter, r *http.Request) {

	var dataPrefix = []byte("data: ")
	var doneSequence = []byte("[DONE]")
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
	req.Header.Set("Authorization", "Bearer "+"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")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// 从另一个Stream API的响应中读取数据并发送给客户端

	lastline := []byte{}
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			line, err := reader.ReadSlice('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				writeError(w, err)
				return
			}
			if len(line) == 1 && line[0] == '\n' {
				continue
			}
			if line[len(line)-2] != '}' {
				lastline = line
				continue
			}
			if len(lastline) > 0 {
				line = append(lastline, line...)
				lastline = []byte{}
			}
			// make sure there isn't any extra whitespace before or after
			line = bytes.TrimSpace(line)
			// the completion API only returns data events
			if !bytes.HasPrefix(line, dataPrefix) {
				continue
			}
			line = bytes.TrimPrefix(line, dataPrefix)

			// the stream is completed when terminated by [DONE]
			if bytes.HasPrefix(line, doneSequence) {
				writeDone(w)
				break
			}

			var res ChatResponse
			err = json.Unmarshal(line, &res)
			if err != nil {
				lastline = line
				// fmt.Println("Error reading response:", err)
				continue
			}
			if res.Choices[0].Delta.Content == "" && res.Choices[0].FinishReason == "stop" {
				return
			}
			writeMessage(w, res.Choices[0].Delta.Content)
			// fmt.Fprintf(w, "%s", res.Choices[0].Delta.Content)
			// w.(http.Flusher).Flush()
		}
	}
}
