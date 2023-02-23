package route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/h-alex2/go-server-test/config"
)

func indexHandler(str string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, str)
	}
}

func NewHttpHandler(client *config.MongoDB, context context.Context) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler("test")).Methods("GET")
	router.HandleFunc("/tasks", CreateTaskHandler(client, context)).Methods("POST")
	router.HandleFunc("/tasks/{id}", GetTaskHandler(client, context)).Methods("GET")
	router.HandleFunc("/tasks/{id}", UpdateTaskHandler(client, context)).Methods("PATCH")
	router.HandleFunc("/tasks/{id}", DeleteTaskHandler(client, context)).Methods("DELETE")

	return router
}
