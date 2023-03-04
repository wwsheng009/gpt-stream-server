package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gpt_stream_server/chatdb"
	"gpt_stream_server/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	req.Header.Set("Authorization", "Bearer "+config.MainConfig.OpenaiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var res ErrorResponse
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error reading response:", err)
			writeError(w, err)
			return
		} else {
			resp.Body.Close()
			writeError(w, fmt.Errorf("%s", res.Error.Message))
			return
		}
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
			line = bytes.TrimSpace(line)
			if len(line) == 1 && line[0] == '\n' {
				continue
			}
			// the completion API only returns data events
			if !bytes.HasPrefix(line, dataPrefix) {
				// lastline = append(lastline, line...)
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
				writeError(w, err)
				continue
			}
			if res.Choices[0].Text == "" && res.Choices[0].FinishReason == "stop" {
				writeDone(w)
				return
			}
			writeMessage(w, res.Choices[0].Text)
		}
	}

}

func processChat(w http.ResponseWriter, r *http.Request, option JsonBody,
	converation chatdb.Conversation) {

	var dataPrefix = []byte("data: ")
	var doneSequence = []byte("[DONE]")

	client := &http.Client{}

	// golang add network proxy
	// 创建一个HTTP客户端

	if config.MainConfig.ProxyServer != "" {
		// Set up the proxy URL
		proxyUrl, _ := url.Parse(config.MainConfig.ProxyServer)
		// Create the HTTP transport with the proxy and other options
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		client.Transport = transport
	}

	var chatdb = chatdb.LoadApiSetting()
	var temp = chatdb.Temperature
	var prompt = option.Prompt
	stopWord := []string{}
	if chatdb.Stop != "" {
		stopWord = append(stopWord, strings.Split(chatdb.Stop, ",")...)
	}
	messages := []Message{{Role: "system", Content: "you are assistant"}}
	if len(converation.Messages) > 0 {
		for _, message := range converation.Messages {
			if message.Prompt != "" {
				messages = append(messages, Message{Role: "user", Content: message.Prompt})
			}
			if message.Completion != "" {
				messages = append(messages, Message{Role: "assistant", Content: message.Completion})
			}
		}
	}
	messages = append(messages, Message{Role: "user", Content: prompt})

	// content = "讲出你的名字"
	request := ChatRequest{
		Model:            chatdb.Model, //"gpt-3.5-turbo",
		Messages:         messages,
		Stream:           true,
		Stop:             stopWord,
		MaxTokens:        &chatdb.MaxTokens,
		Temperature:      &temp,
		TopP:             &chatdb.TopP,
		PresencePenalty:  chatdb.PresencePenalty,
		FrequencyPenalty: chatdb.FrequencyPenalty,
		N:                &chatdb.N,
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
	req.Header.Set("Authorization", "Bearer "+chatdb.ApiToken) //config.MainConfig.OpenaiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var res ErrorResponse
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Println("Error reading response:", err)
			writeError(w, err)
			return
		} else {
			resp.Body.Close()
			writeError(w, fmt.Errorf("%s", res.Error.Message))
			return
		}
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// 从另一个Stream API的响应中读取数据并发送给客户端

	// lastline := []byte{}
	answser := ""
	writeConversationId(w, converation.ConversationId)

	for {
		select {
		case <-r.Context().Done():
			SaveText(converation.Id, prompt, answser)
			return
		default:
			line, err := reader.ReadSlice('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				writeError(w, err)
				return
			}
			// make sure there isn't any extra whitespace before or after
			line = bytes.TrimSpace(line)
			if len(line) == 1 && line[0] == '\n' {
				continue
			}
			if !bytes.HasPrefix(line, dataPrefix) {
				// lastline = append(lastline, line...)
				continue
			}
			line = bytes.TrimPrefix(line, dataPrefix)

			// the stream is completed when terminated by [DONE]
			if bytes.HasPrefix(line, doneSequence) {
				writeDone(w)
				SaveText(converation.Id, prompt, answser)
				break
			}

			var res ChatResponse
			err = json.Unmarshal(line, &res)
			if err != nil {
				// lastline = line
				// fmt.Println("Error reading response:", err)
				writeError(w, err)
				continue
			}
			if len(res.Choices) > 0 {
				text := res.Choices[0].Delta.Content
				if text == "" && res.Choices[0].FinishReason == "stop" {
					SaveText(converation.Id, prompt, answser)
					return
				}
				writeMessage(w, text)
				answser += text
			} else {
				if res.Error.Message != "" {
					writeError(w, errors.New(res.Error.Message))
				}
			}
		}
	}

}

func writeConversationId(w http.ResponseWriter, s string) {
	fmt.Fprintf(w, `{"conversationId":"%s"}`, s)
	w.Write([]byte("\n"))
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}

func SaveText(conversationId int32, prompt, text string) {
	if len(text) == 0 {
		return
	}
	// println(text)
	chatdb.CreateNewMessage(conversationId, prompt, text)
}
