package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// const (
// 	ApiServer          = "API_SERVER"
// 	ApiServerAccessKey = "API_SERVER_ACESS_KEY"
// 	HttpPort           = "HTTP_PORT"
// 	HttpHost           = "HTTP_HOST"
// 	StaticFolder       = "STATIC_FOLDER"
// 	OpenaiKey          = "OPENAI_KEY"
// 	ProxyServer        = "PROXY_SERVER"
// )

type Config struct {
	// api server to get config
	ApiServer string `json:"api_server"`
	//access key
	ApiServerAccessKey string `json:"api_server_access_key"`
	// server port
	HttpPort string `json:"http_port"`
	//server host
	HttpHost string `json:"http_host"`
	//static folder location
	StaticFolder string `json:"static_folder"`
	//openai key
	OpenaiKey string `json:"openai_key"`
	//use network proxy
	HttpsProxy string `json:"proxy_server"`
	Storage    string `json:"storeage"`
}

var MainConfig = Config{}

// func NewConfig() *Config {
// 	MainConfig = Config{
// 		ApiServer:    "https://api.example.com",
// 		HttpPort:     "8080",
// 		HttpHost:     "localhost",
// 		StaticFolder: "./static",
// 		OpenaiKey:    "",
// 		ProxyServer:  "",
// 	}
// 	LoadConfigFromEnv(&MainConfig)
// 	return &MainConfig
// }

// InitConfig is a function that initializes and loads the configuration data from .env file.
func InitConfig() {
	// Load the .env file
	err := godotenv.Load(".env")

	// Load the .env file
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	// Load the configuration data from the .env file
	LoadConfigFromEnv(&MainConfig)
}

func LoadConfigFromEnv(cfg *Config) error {
	var ok bool
	if cfg.ApiServer, ok = os.LookupEnv("API_SERVER"); !ok {
		return errors.New("API_SERVER is required")
	}
	if cfg.ApiServerAccessKey, ok = os.LookupEnv("API_SERVER_ACESS_KEY"); !ok {
		return errors.New("API_SERVER_ACESS_KEY is required")
	}

	cfg.HttpPort, _ = os.LookupEnv("HTTP_PORT")

	cfg.HttpHost, _ = os.LookupEnv("HTTP_HOST")

	cfg.StaticFolder, _ = os.LookupEnv("STATIC_FOLDER")
	cfg.OpenaiKey, _ = os.LookupEnv("OPENAI_KEY")
	cfg.HttpsProxy, _ = os.LookupEnv("HTTPS_PROXY")
	if cfg.HttpsProxy == "" {
		cfg.HttpsProxy, _ = os.LookupEnv("ALL_PROXY")
	}

	cfg.Storage, _ = os.LookupEnv("STORAGE")
	return nil
}
