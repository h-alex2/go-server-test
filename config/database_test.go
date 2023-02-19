package config

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetClientMongoDB(t *testing.T) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_URL := os.Getenv("DB_URI")
	t.Run("DB에 연결돼야 한다.", func(t *testing.T) {
		db, err := GetMongoDBClient(DB_URL)
		assert.NoError(t, err)
		assert.NotNil(t, db.client)
	})

	t.Run("주소가 옳지 않으면 에러를 내뱉어야 한다.", func(t *testing.T) {
		inCorrectUri := "mongodb://localhost:99999"
		db, err := GetMongoDBClient(inCorrectUri)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
