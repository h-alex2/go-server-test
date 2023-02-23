package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/h-alex2/go-server-test/config"
	"github.com/h-alex2/go-server-test/route"
	"github.com/h-alex2/go-server-test/util"
)

func main() {
	DB_URL := util.GetEnv("DB_URI")
	DB_NAME := util.GetEnv("DB_NAME")
	BASE_URI := util.GetEnv("BASE_URI")

	client, err, context, cancel := config.GetMongoDBClient(DB_URL, DB_NAME)

	if err != nil {
		panic(err)
	}

	defer func() {
		client.CloseMongoDB()
		cancel()
	}()

	fmt.Println("Success Connect DB")

	srv := &http.Server{
		Addr:         BASE_URI,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      route.NewHttpHandler(client, context),
	}

	func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()
}
