package gpt

import (
	"context"
	"errors"
	"fmt"
	"gpt_stream_server/chatdb"
	"gpt_stream_server/config"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func StreamHandler(c *gin.Context) {
	// 设置响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	// 读取请求体
	var jsonBody = new(JsonBody)
	if c.Request.Method == "POST" {
		// 解析JSON
		err := c.ShouldBindJSON(jsonBody)
		if err != nil {
			// 处理错误
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		//new conversation
		var conversation chatdb.Conversation //{}
		if jsonBody.Option.ConversationId == "" {
			conversation = chatdb.CreateNewconversation(jsonBody.Prompt)
		} else {
			conversation = chatdb.FindConversationById(jsonBody.Option.ConversationId)
		}
		processChat(c.Writer, c.Request, *jsonBody, conversation)

	}
	// processComplete(c.Writer, c.Request)

	// 处理请求
}

type JsonBody struct {
	Prompt string `json:"prompt"`
	Option struct {
		ConversationId  string `json:"conversationId"`
		ParentMessageId string `json:"parentMessageId"`
	} `json:"options"`
}

func gpt3client(c *gin.Context) {
	var config = gogpt.DefaultConfig(config.MainConfig.OpenaiKey)

	var client = gogpt.NewClientWithConfig(config)

	ctx := context.Background()

	var jsonBody = new(JsonBody)
	if c.Request.Method == "POST" {
		// 解析JSON
		err := c.ShouldBindJSON(jsonBody)
		if err != nil {
			// 处理错误
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

	}

	request := gogpt.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo",
		Messages:  []gogpt.ChatCompletionMessage{{Role: "system", Content: "you are ai asistant"}, {Role: "user", Content: jsonBody.Prompt}},
		MaxTokens: 2048,
	}
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create request"})
		return
	}
	defer stream.Close()
	for {

		//2023-3-3 这个库存在问题，如果调用接口异常出错，没有判断，也没有返回错误
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}
		writeMessage(c.Writer, response.Choices[0].Delta.Content)
		// fmt.Printf("Stream response: %v\n", response)
	}
}