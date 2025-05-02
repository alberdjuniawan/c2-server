package config

import (
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
)

var JWTSecret string
var DBPath string

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return err
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	DBPath = os.Getenv("DB_PATH")

	if JWTSecret == "" || DBPath == "" {
		log.Fatal("Environment variables JWT_SECRET or DB_PATH are missing")
		return fmt.Errorf("missing required environment variables")
	}

	return nil
}
