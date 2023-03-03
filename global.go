package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	api_proxy string
	http_port string
	http_host string
}

var config = Config{}

func Init_config() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbName := os.Getenv("DB_NAME")
	// dbUser := os.Getenv("DB_USER")
	config.api_proxy = os.Getenv("API_SERVER")
	config.http_port = os.Getenv("HTTP_PORT")
	config.http_host = os.Getenv("HTTP_HOST")
	if len(config.http_port) == 0 {
		config.http_port = "8080"
	}

}
