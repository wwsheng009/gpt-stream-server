package chat

import (
	"encoding/json"
	"fmt"
	"gpt_stream_server/config"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ConfigState struct {
	ApiModel     string `json:"apiModel"`
	ReverseProxy string `json:"reverseProxy"`
	SocksProxy   string `json:"socksProxy"`
	HttpsProxy   string `json:"httpsProxy"`
	Balance      string `json:"balance"`
	TimeoutMs    int64  `json:"timeoutMs"`
}
type ModelConfig struct {
	Status string      `json:"status"`
	Data   ConfigState `json:"data"`
}

func CheckAuth(c *gin.Context) error {
	authSecretKey := os.Getenv("AUTH_SECRET_KEY")
	if len(authSecretKey) > 0 {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || strings.TrimSpace(strings.ReplaceAll(authorization, "Bearer ", "")) != authSecretKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "Unauthorized",
				"message": "Error: 无访问权限 | No access rights",
				"data":    nil,
			})
			return fmt.Errorf("error: 无访问权限")
		}
	}
	return nil
}
func Auth(c *gin.Context) {
	authSecretKey := os.Getenv("AUTH_SECRET_KEY")
	if len(authSecretKey) > 0 {
		authorization := c.GetHeader("Authorization")
		if authorization == "" || strings.TrimSpace(strings.ReplaceAll(authorization, "Bearer ", "")) != authSecretKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "Unauthorized",
				"message": "Error: 无访问权限 | No access rights",
				"data":    nil,
			})
			return
		}
	}
	c.Next()
}

func ConfigHandler(c *gin.Context) {
	response, err := chatConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func getTimeOutMs() int64 {
	timeoutMsEnv := os.Getenv("TIMEOUT_MS")
	timeoutMs, err := strconv.Atoi(timeoutMsEnv)
	if err != nil {
		timeoutMs = 30 * 1000
	}
	timeoutDuration := time.Duration(timeoutMs) * time.Millisecond
	return timeoutDuration.Microseconds()
}

func chatConfig() (interface{}, error) {
	balance, err := fetchBalance()
	if err != nil {
		return nil, err
	}
	reverseProxy := os.Getenv("API_REVERSE_PROXY")
	if reverseProxy == "" {
		reverseProxy = "-"
	}
	httpsProxy := config.MainConfig.HttpsProxy

	if httpsProxy == "" {
		httpsProxy = "-"
	}
	socksProxy := "-"
	if host := os.Getenv("SOCKS_PROXY_HOST"); host != "" {
		port := os.Getenv("SOCKS_PROXY_PORT")
		if port != "" {
			socksProxy = host + ":" + port
		}
	}
	return ModelConfig{
		Status: "Success",
		Data: ConfigState{
			ApiModel: CurrentModel(), ReverseProxy: reverseProxy, SocksProxy: socksProxy,
			HttpsProxy: httpsProxy, Balance: balance, TimeoutMs: getTimeOutMs(),
		},
	}, nil
}

func fetchBalance() (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if !isNotEmptyString(apiKey) {
		return "-", nil
	}
	apiBaseURL := os.Getenv("OPENAI_API_BASE_URL")
	if apiBaseURL == "" {
		apiBaseURL = "https://api.openai.com"
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiBaseURL+"/dashboard/billing/credit_grants", nil)
	if err != nil {
		return "-", nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)

	if err != nil {
		return "-", nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "-", nil
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "-", nil
	}
	balance, ok := data["total_available"].(float64)
	if !ok {
		return "-", nil
	}
	return fmt.Sprintf("%.3f", balance), nil
}
