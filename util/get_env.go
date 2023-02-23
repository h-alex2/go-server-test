package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(envKey string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(envKey)
}
