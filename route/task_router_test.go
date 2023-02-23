package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/h-alex2/go-server-test/config"
	"github.com/h-alex2/go-server-test/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTaskHandler(t *testing.T) {
	assert := assert.New(t)
	var dbUri string = "mongodb://localhost:27017"

	client, err, ctx, cancel := config.GetMongoDBClient(dbUri, "testing_db")

	if err != nil {
		t.Fatal(err)
	}

	collection := client.GetCollection("tasks")
	t.Run("GET, id에 맞는 task를 전달해야 한다.", func(t *testing.T) {
		filter := bson.D{{"name", "test"}, {"total_time", 0}}
		id, _ := collection.InsertOne(ctx, filter)
		var createdTestDocumentId primitive.ObjectID = id.InsertedID.(primitive.ObjectID)

		var createdTestDocument bson.M
		_ = client.GetCollection("tasks").FindOne(ctx, filter).Decode(&createdTestDocument)

		createdTestDocumentJson, _ := json.Marshal(createdTestDocument)

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/tasks/"+createdTestDocumentId.Hex(), nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)
		responseJsonData, _ := ioutil.ReadAll(res.Body)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(createdTestDocumentJson, responseJsonData)

		defer func() {
			collection.DeleteOne(ctx, bson.M{"_id": createdTestDocumentId})
			client.CloseMongoDB()
			cancel()
		}()
	})

	t.Run("GET, id에 맞는 task가 없으면 400 에러를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/tasks/inValidId", nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)

		assert.Equal(http.StatusBadRequest, res.Code)
	})
}

func TestCreateTaskHandler(t *testing.T) {
	assert := assert.New(t)
	var dbUri string = "mongodb://localhost:27017"

	client, err, ctx, cancel := config.GetMongoDBClient(dbUri, "testing_db")

	if err != nil {
		t.Fatal(err)
	}

	collection := client.GetCollection("tasks")
	t.Run("POST, 생성하고, 생성한 문서의 id를 전달해야 한다.", func(t *testing.T) {
		ts := httptest.NewServer(NewHttpHandler(client, ctx))

		res, _ := http.Post(ts.URL+"/tasks", "application/json", strings.NewReader(`{"name": "test", "total_time": 0}`))

		var result1 bson.M
		var result2 bson.M

		err := collection.FindOne(ctx, bson.D{{"name", "test"}}).Decode(&result1)
		responseDocumentId, _ := ioutil.ReadAll(res.Body)
		id, _ := primitive.ObjectIDFromHex(string(responseDocumentId))

		idErr := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result2)

		assert.NoError(err)
		assert.NoError(idErr)
		assert.Equal(http.StatusCreated, res.StatusCode)

		defer collection.DeleteMany(ctx, bson.M{"name": "test"})
	})

	t.Run("POST, request의 json에서 total_time 필드가 비어있어도 기본값으로 설정되어야 한다.", func(t *testing.T) {
		ts := httptest.NewServer(NewHttpHandler(client, ctx))

		res, _ := http.Post(ts.URL+"/tasks", "application/json", strings.NewReader(`{"name": "test2"}`))

		var result1 bson.M
		var result2 bson.M

		err := collection.FindOne(ctx, bson.D{{"name", "test2"}}).Decode(&result1)
		responseDocumentId, _ := ioutil.ReadAll(res.Body)
		id, _ := primitive.ObjectIDFromHex(string(responseDocumentId))
		idErr := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result2)

		assert.NoError(err)
		assert.NoError(idErr)
		assert.Equal(result2["total_time"], int32(0))
		assert.Equal(http.StatusCreated, res.StatusCode)

		defer collection.DeleteMany(ctx, bson.M{"name": "test2"})
	})

	t.Run("POST, request의 json에서 name 필드가 비어있으면 400 에러를 전달해야 한다.", func(t *testing.T) {
		ts := httptest.NewServer(NewHttpHandler(client, ctx))

		res, _ := http.Post(ts.URL+"/tasks", "application/json", strings.NewReader(""))

		assert.Equal(http.StatusBadRequest, res.StatusCode)
	})

	defer func() {
		client.CloseMongoDB()
		cancel()
	}()
}

func TestUpdateTaskHandler(t *testing.T) {
	assert := assert.New(t)
	var dbUri string = "mongodb://localhost:27017"

	client, err, ctx, cancel := config.GetMongoDBClient(dbUri, "testing_db")

	if err != nil {
		t.Fatal(err)
	}

	collection := client.GetCollection("tasks")
	filter := bson.D{{"name", "test"}}
	id, _ := collection.InsertOne(ctx, filter)
	var createdTestDocumentId primitive.ObjectID = id.InsertedID.(primitive.ObjectID)

	t.Run("PATCH, 수정된 문서를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PATCH", "/tasks/"+createdTestDocumentId.Hex(), strings.NewReader(`{"name": "test3"}`))
		req.Header.Set("Content-Type", "application/json")

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)

		responseDocument, _ := ioutil.ReadAll(res.Body)

		var updatedTask model.Task
		_ = json.Unmarshal(responseDocument, &updatedTask)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal("test3", updatedTask.Name)

		defer collection.DeleteMany(ctx, bson.M{"name": "test3"})
	})

	t.Run("PATCH, request의 필드가 비어있으면 400 에러를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PATCH", "/tasks/"+createdTestDocumentId.Hex(), nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)
		assert.Equal(http.StatusBadRequest, res.Code)
	})

	t.Run("PATCH, id에 맞는 task가 없으면 400 에러를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PATCH", "/tasks/"+"invalidId", nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)
		assert.Equal(http.StatusBadRequest, res.Code)
	})

	defer func() {
		defer collection.DeleteMany(ctx, bson.M{"name": "test"})
		client.CloseMongoDB()
		cancel()
	}()
}

func TestDeleteTaskHandler(t *testing.T) {
	assert := assert.New(t)
	var dbUri string = "mongodb://localhost:27017"

	client, err, ctx, cancel := config.GetMongoDBClient(dbUri, "testing_db")

	if err != nil {
		t.Fatal(err)
	}

	collection := client.GetCollection("tasks")
	filter := bson.D{{"name", "deleteTest"}}
	id, _ := collection.InsertOne(ctx, filter)
	var createdTestDocumentId primitive.ObjectID = id.InsertedID.(primitive.ObjectID)

	t.Run("DELETE, 204 코드를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/tasks/"+createdTestDocumentId.Hex(), nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)
		assert.Equal(http.StatusNoContent, res.Code)
	})

	t.Run("DELETE, id에 맞는 문서를 삭제해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/tasks/"+createdTestDocumentId.Hex(), nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)

		result := bson.M{}
		err := client.GetCollection("tasks").FindOne(ctx, bson.D{{"_id", createdTestDocumentId}}).Decode(&result)

		assert.Equal(http.StatusNoContent, res.Code)
		assert.Error(err)
	})

	t.Run("DELETE, id에 맞는 task가 없으면 400 에러를 전달해야 한다.", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/tasks/"+"invalidId", nil)

		mux := NewHttpHandler(client, ctx)
		mux.ServeHTTP(res, req)
		assert.Equal(http.StatusBadRequest, res.Code)
	})

	defer func() {
		defer collection.DeleteMany(ctx, bson.M{"name": "deleteTest"})
		client.CloseMongoDB()
		cancel()
	}()
}
