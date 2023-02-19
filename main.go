package main

import (
	"fmt"
	"log"
	"os"

	"github.com/h-alex2/go-server-test/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_URL := os.Getenv("DB_URI")

	client, err := config.GetMongoDBClient(DB_URL)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.CloseMongoDB(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Success Connect DB")
}
