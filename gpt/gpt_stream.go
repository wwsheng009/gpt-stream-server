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
	"net/http"
	"net/url"
	"strings"
	"time"
)

func processRequest(w http.ResponseWriter, r *http.Request, option RequestBody,
	converation *chatdb.Conversation, setting *chatdb.ApiSetting) {

	prompt := option.Prompt
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

	// var chatdb = chatdb.LoadApiSetting()

	isChat := false
	url := "https://api.openai.com/v1/completions"
	if strings.Contains(setting.Model, "gpt-3.5-turbo") {
		isChat = true
		url = "https://api.openai.com/v1/chat/completions"
	}

	start := time.Now()
	buf, err := getRequestBuf(isChat, prompt, setting, converation)
	if err != nil {
		writeError(w, err)
		return
	}
	// 发送另一个Stream API的请求
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		writeError(w, err)
		return
	}
	req.Header.Set("Accept", "text/event-stream; charset=utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+setting.ApiToken) //config.MainConfig.OpenaiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		writeError(w, err)
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
			SaveText(converation, prompt, answser, start)
			return
		default:
			line, err := reader.ReadSlice('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Println("Error reading response:", err)
				writeError(w, err)
				return
			}
			// make sure there isn't any extra whitespace before or after
			line = bytes.TrimSpace(line)
			// if len(line) == 1 && line[0] == '\n' {
			// 	continue
			// }
			if !bytes.HasPrefix(line, dataPrefix) {
				// lastline = append(lastline, line...)
				continue
			}
			line = bytes.TrimPrefix(line, dataPrefix)

			// the stream is completed when terminated by [DONE]
			if bytes.HasPrefix(line, doneSequence) {
				// writeDone(w)
				SaveText(converation, prompt, answser, start)
				break
			}
			text, err := processLine(isChat, line)
			if err != nil {
				writeError(w, err)
				continue
			}
			writeMessage(w, text)
			answser += text
		}
	}

}
func getRequestBuf(isChat bool, prompt string, chatdb *chatdb.ApiSetting, converation *chatdb.Conversation) (bytes.Buffer, error) {
	if isChat {
		return getChatBuf(prompt, chatdb, converation)
	} else {
		return getCompletionBuf(prompt, chatdb, converation)
	}
}
func getCompletionBuf(prompt string, chatdb *chatdb.ApiSetting, converation *chatdb.Conversation) (bytes.Buffer, error) {
	var temp = chatdb.Temperature
	stopWord := []string{}
	if chatdb.Stop != "" {
		stopWord = append(stopWord, strings.Split(chatdb.Stop, ",")...)
	}
	messages := "" // "提示:你叫" + chatGptName + "。\n"

	if len(chatdb.AiNickname) > 0 {
		messages += "提示:你叫" + chatdb.AiNickname + "。\n"
	}
	if len(converation.Messages) > 0 {
		for _, message := range converation.Messages {

			if message.Prompt != "" {
				if len(chatdb.UserNickname) > 0 {
					messages += chatdb.UserNickname + ": "
				}
				messages += message.Prompt + "\n\n"
			}
			if message.Completion != "" {
				if len(chatdb.AiNickname) > 0 {
					messages += chatdb.AiNickname + ": "
				}
				messages += message.Completion + "\n\n"
			}
		}
	}
	if len(chatdb.UserNickname) > 0 {
		messages += chatdb.UserNickname + ": "
	}
	messages += prompt + "\n\n"
	if len(chatdb.AiNickname) > 0 {
		messages += chatdb.AiNickname + ": "
	}

	// content = "讲出你的名字"
	request := CompletionRequest{
		Model:            chatdb.Model, //"gpt-3.5-turbo",
		Prompt:           messages,
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
		return buf, err
	}
	return buf, nil

}
func getChatBuf(prompt string, chatdb *chatdb.ApiSetting, converation *chatdb.Conversation) (bytes.Buffer, error) {
	var temp = chatdb.Temperature
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
		return buf, err
	}
	return buf, nil
}
func processLine(isChat bool, line []byte) (string, error) {
	if isChat {
		return processChatLine(line)
	} else {
		return processCompletionLine(line)
	}

}
func processChatLine(line []byte) (string, error) {
	text := ""
	var res ChatResponse
	err := json.Unmarshal(line, &res)
	if err != nil {
		// lastline = line
		// fmt.Println("Error reading response:", err)
		return "", err
	}
	if len(res.Choices) > 0 {
		text = res.Choices[0].Delta.Content
	} else {
		if res.Error.Message != "" {
			return "", errors.New(res.Error.Message)
		}
	}
	return text, nil
}

func processCompletionLine(line []byte) (string, error) {
	text := ""
	var res CompletionResponse
	err := json.Unmarshal(line, &res)
	if err != nil {
		// lastline = line
		// fmt.Println("Error reading response:", err)
		return "", err
	}
	if len(res.Choices) > 0 {
		text = res.Choices[0].Text
	} else {
		if res.Error.Message != "" {
			return "", errors.New(res.Error.Message)
		}
	}
	return text, nil
}

func writeConversationId(w http.ResponseWriter, s string) {
	fmt.Fprintf(w, `{"conversationId":"%s"}`, s)
	w.Write([]byte("\n"))
	w.WriteHeader(http.StatusOK)
	w.(http.Flusher).Flush()
}

func SaveText(converation *chatdb.Conversation, prompt, text string, start time.Time) {
	if len(text) == 0 {
		return
	}

	end := time.Now()
	delta := end.Sub(start)

	// println(text)
	chatdb.GetDefaultConversation().CreateNewMessage(converation, prompt, text, delta.Seconds())
}
