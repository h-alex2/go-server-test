package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/h-alex2/go-server-test/config"
	"github.com/h-alex2/go-server-test/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTaskHandler(client *config.MongoDB, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, idErr := primitive.ObjectIDFromHex(vars["id"])

		if idErr != nil {
			fmt.Println(idErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filter := bson.D{{"_id", id}}
		var result bson.M

		err := client.GetCollection("tasks").FindOne(context, filter).Decode(&result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data, _ := json.Marshal(result)

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(data))
	}
}

func CreateTaskHandler(client *config.MongoDB, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := new(model.Task)
		err := json.NewDecoder(r.Body).Decode(task)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		task.CreatedAt = time.Now()

		id, _ := client.GetCollection("tasks").InsertOne(context, task)
		var createdDocumentId primitive.ObjectID = id.InsertedID.(primitive.ObjectID)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, createdDocumentId.Hex())
	}
}

func UpdateTaskHandler(client *config.MongoDB, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, idErr := primitive.ObjectIDFromHex(vars["id"])

		var filter bson.M = bson.M{"_id": id}

		var reqData bson.M = bson.M{}

		err := json.NewDecoder(r.Body).Decode(&reqData)
		var update bson.M = bson.M{"$set": reqData}

		if err != nil || idErr != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}

		var result bson.M
		_ = client.GetCollection("tasks").FindOneAndUpdate(context, filter, update, &opt).Decode(&result)

		data, _ := json.Marshal(result)

		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json")
		fmt.Fprint(w, string(data))
	}
}

func DeleteTaskHandler(client *config.MongoDB, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, idErr := primitive.ObjectIDFromHex(vars["id"])

		if idErr != nil {
			fmt.Println(id, idErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filter := bson.D{{"_id", id}}

		_, err := client.GetCollection("tasks").DeleteOne(context, filter)

    if err != nil {
			fmt.Println(id, idErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
