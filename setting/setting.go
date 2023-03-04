package setting

import (
	"encoding/json"
	"gpt_stream_server/yao"
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

func LoadApiSetting() *ApiSetting {

	obj, err := yao.YaoProcess("scripts.ai.chatgpt.GetSetting")
	if err != nil {
		panic(err.Error())
	}

	d, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	setting := new(ApiSetting)
	err = json.Unmarshal(d, setting)
	if err != nil {
		panic(err.Error())
	}

	// fmt.Println(setting.ApiToken)
	return setting
}
