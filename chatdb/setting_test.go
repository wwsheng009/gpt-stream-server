package chatdb

import (
	"gpt_stream_server/config"
	"testing"
)

func TestLoadApiSetting(t *testing.T) {
	config.InitConfig()
	LoadApiSetting()

}
