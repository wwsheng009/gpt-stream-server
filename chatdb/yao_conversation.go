package chatdb

import (
	"encoding/json"
	"gpt_stream_server/yao"
)

type YaoConversation struct {
	Setting *ApiSetting
}

func (y *YaoConversation) LoadApiSetting() (*ApiSetting, error) {

	obj, err := yao.YaoProcess("scripts.ai.chatgpt.GetSetting")
	if err != nil {
		panic(err.Error())
	}
	setting := ApiSetting{}
	err = ConvertData(obj, &setting)
	if err != nil {
		panic(err.Error())
	}

	if setting.MaxSendLines > 20 {
		setting.MaxSendLines = 20
	}
	if setting.MaxSendLines < 1 {
		setting.MaxSendLines = 1
	}

	y.Setting = &setting
	return &setting, nil
}

func (y *YaoConversation) CreateNewconversation(title string) (*Conversation, error) {
	obj, err := yao.YaoProcess("scripts.chat.conversation.NewConversation", title)
	if err != nil {
		return nil, err
	}
	conversation := Conversation{}
	err = ConvertData(obj, &conversation)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (y *YaoConversation) FindConversationById(conversationId string) (*Conversation, error) {
	obj, err := yao.YaoProcess("scripts.chat.conversation.FindConversationById", conversationId)
	if err != nil {
		return nil, err
	}
	conversation := Conversation{}
	err = ConvertData(obj, &conversation)
	if err != nil {
		return nil, err
	}

	conversation.Messages = GetLastLines(conversation.Messages, y.Setting.MaxSendLines)

	return &conversation, nil
}

func (y *YaoConversation) CreateNewMessage(converation *Conversation, prompt string, answer string, seconds float64) error {

	request := map[string]interface{}{
		"conversationId": converation.Id,
		"prompt":         prompt,
		"answer":         answer,
		"seconds":        seconds,
	}
	_, err := yao.YaoProcess("scripts.chat.conversation.NewMessageApi", request)
	if err != nil {
		// panic(err.Error())
		return err
	}
	return nil
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
