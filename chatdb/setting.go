package chatdb

import (
	"fmt"
	"gpt_stream_server/yao"
	"os"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const expiryDateLayout = "2006-01-02 15:04:05"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(expiryDateLayout, s)
	return
}

type ApiSetting struct {
	AccessCount      int         `json:"access_count"`
	AiNickname       string      `json:"ai_nickname"`
	ApiToken         string      `json:"api_token"`
	CreatedAt        CustomTime  `json:"created_at"`
	Default          bool        `json:"default"`
	DeletedAt        interface{} `json:"deleted_at"` // can be null or a time.Time value
	Description      string      `json:"description"`
	FrequencyPenalty float32     `json:"frequency_penalty"`
	Id               int         `json:"id"`
	MaxSendLines     int         `json:"max_send_lines"`
	MaxTokens        int         `json:"max_tokens"`
	Model            string      `json:"model"`
	PresencePenalty  float32     `json:"presence_penalty"`
	Stop             string      `json:"stop"`
	Temperature      float32     `json:"temperature"`
	TopP             float32     `json:"top_p"`
	UpdatedAt        CustomTime  `json:"updated_at"`
	UserNickname     string      `json:"user_nickname"`
	N                int         `json:"n"`
}

func LoadApiSetting() ApiSetting {

	obj, err := yao.YaoProcess("scripts.ai.chatgpt.GetSetting")
	if err != nil {
		panic(err.Error())
	}
	setting := ApiSetting{}
	err = ConvertData(obj, &setting)
	if err != nil {
		panic(err.Error())
	}
	return setting
}

func LoadLocalApiSetting() (*ApiSetting, error) {
	if _, err := os.Stat("path/to/file"); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return nil, err
	}
	data, err := os.ReadFile("./gpt.config.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	setting := ApiSetting{}
	err = ConvertData(data, &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}
