package chatdb

import (
	"encoding/json"
	"gpt_stream_server/yao"
)

type Messages struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}
type Conversation struct {
	ConversationId string     `json:"uuid"`
	Id             int32      `json:"id"`
	Title          string     `json:"title"`
	Messages       []Messages `json:"messages"`
}

func CreateNewconversation(prompt string) Conversation {
	obj, err := yao.YaoProcess("scripts.chat.conversation.NewConversation", prompt)
	if err != nil {
		panic(err.Error())
	}
	conversation := Conversation{}
	err = ConvertData(obj, &conversation)
	if err != nil {
		panic(err.Error())
	}
	return conversation
}

func FindConversationById(conversationId string) Conversation {
	obj, err := yao.YaoProcess("scripts.chat.conversation.FindConversationById", conversationId)
	if err != nil {
		panic(err.Error())
	}
	conversation := Conversation{}
	err = ConvertData(obj, &conversation)
	if err != nil {
		panic(err.Error())
	}
	return conversation
}

func CreateNewMessage(conversationId int32, prompt string, answer string) {

	request := map[string]interface{}{
		"conversationId": conversationId,
		"prompt":         prompt,
		"answer":         answer,
	}
	_, err := yao.YaoProcess("scripts.chat.conversation.NewMessageApi", request)
	if err != nil {
		panic(err.Error())
	}
}

func ConvertData(obj interface{}, ref interface{}) error {
	d, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, ref)
	if err != nil {
		return err
	}
	return nil
}
