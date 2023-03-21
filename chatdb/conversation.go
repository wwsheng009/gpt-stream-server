package chatdb

import (
	"gpt_stream_server/config"
	"strings"
)

type IConversation interface {
	LoadApiSetting() (*ApiSetting, error)
	CreateNewconversation(title string) (*Conversation, error)
	FindConversationById(conversationId string) (*Conversation, error)
	CreateNewMessage(converation *Conversation, prompt string, answer string, seconds float64) error
}

type ConvMessage struct {
	Prompt     string  `json:"prompt"`
	Completion string  `json:"completion"`
	Seconds    float64 `json:"seconds"`
}
type Conversation struct {
	ConversationId string        `json:"uuid"`
	Id             int32         `json:"id"`
	Title          string        `json:"title"`
	Messages       []ConvMessage `json:"messages"`
}

func GetDefaultConversation() IConversation {
	storage := config.MainConfig.Storage
	if strings.ToLower(storage) == "local" {
		return &LocalConversation{}
	}
	return &YaoConversation{}
}

func IsLocal() bool {
	storage := config.MainConfig.Storage
	return strings.ToLower(storage) == "local"
}
