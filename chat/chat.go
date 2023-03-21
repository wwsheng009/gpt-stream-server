package chat

import (
	"gpt_stream_server/chatdb"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type SessionData struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Session(c *gin.Context) {
	AUTH_SECRET_KEY := os.Getenv("AUTH_SECRET_KEY")
	hasAuth := isNotEmptyString(AUTH_SECRET_KEY)
	if hasAuth {
		c.JSON(http.StatusOK, SessionData{
			Status:  "Success",
			Message: "",
			Data: map[string]interface{}{
				"auth":  true,
				"model": CurrentModel(),
			},
		})
	} else {
		c.JSON(http.StatusOK, SessionData{
			Status:  "Success",
			Message: "",
			Data: map[string]interface{}{
				"auth":  false,
				"model": CurrentModel(),
			},
		})
	}
}
func isNotEmptyString(str string) bool {
	return str != ""
}

func CurrentModel() string {
	conv := chatdb.GetDefaultConversation()
	setting, err := conv.LoadApiSetting()
	if err != nil {
		return ""
	}

	return setting.Model
}
