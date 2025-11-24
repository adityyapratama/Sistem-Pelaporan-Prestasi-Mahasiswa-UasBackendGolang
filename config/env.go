package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found (menggunakan environment system)")
	}
}

func GetPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}