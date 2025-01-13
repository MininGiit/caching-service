package http

import (
	"bytes"
	"cachingService/internal/infrastructure/logger"
	mock "cachingService/internal/mocks/cache"
	"cachingService/internal/usecase"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGet(t *testing.T) {
	resValue := 7.0
	resKey := "test"
	var resErr error = nil
	cache := &mock.MockCache{
		Value: resValue,
		Key:   resKey,
		Err:   resErr,
	}
	expectedResponse := NewResponseItem(resKey, resValue, time.Now())
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	req, _ := http.NewRequest("GET", "/api/lru/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var response ResponseItem
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, http.StatusOK, recorder.Code, "status code should be 200")
	assert.Equal(t, expectedResponse.Value, response.Value, "Value should match")
}

func TestHandlerGetNotFound(t *testing.T) {
	cache := &mock.MockCache{}
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	request, _ := http.NewRequest("GET", "/api/lru/te", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, "status code should be 404")
}

func TestHandlerPost(t *testing.T) {
	cache := &mock.MockCache{}
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	requestBody := []byte(`{"key": "validKey", "value": [1, 3, 4], "ttlSeconds": 1}`) // Correct json body
	request, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody))   // Create post request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, "Status code should be 400")
}

func TestHandlerPostBad(t *testing.T) {
	cache := &mock.MockCache{}
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	requestBody := []byte(`{"key": "validKey", "value": "test", "ttlSeconds": 1}`)  // Correct json body
	request, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody)) // Create post request
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusCreated, recorder.Code, "Status code should be 201")
}

func TestHandlerDelete(t *testing.T) {
	resKey := "test"
	cache := &mock.MockCache{Key: resKey}
	uc := usecase.New(cache)
	logger := logger.New("error")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	requestBody := []byte(`{"key": "validKey", "value": [1, 3, 4], "ttlSeconds": 1}`)
	request, _ := http.NewRequest("DELETE", "/api/lru/test", bytes.NewReader(requestBody))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNoContent, recorder.Code, "Status code should be 204")
}

func TestHandlerDeleteNotFound(t *testing.T) {
	resKey := "test"
	cache := &mock.MockCache{Key: resKey}
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	requestBody := []byte(`{"key": "validKey", "value": [1, 3, 4], "ttlSeconds": 1}`)
	request, _ := http.NewRequest("DELETE", "/api/lru/tt", bytes.NewReader(requestBody))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, "Status code should be 404")
}

func TestHandlerDeleteAll(t *testing.T) {
	cache := &mock.MockCache{Size: 0}
	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	request, _ := http.NewRequest("DELETE", "/api/lru", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNoContent, recorder.Code, "Status code should be 204")
}

func TestHandlerGetAll(t *testing.T) {
	expectedKeys := []string{"key1", "key2", "key3"}
	expectedValues := []interface{}{1, "2", 3.0}
	cache := &mock.MockCache{
		Size:   10,
		Keyes:  expectedKeys,
		Values: expectedValues,
	}
	expectedResponse := NewResponseItems(expectedKeys, expectedValues)

	uc := usecase.New(cache)
	logger := logger.New("ERROR")
	handler := NewHandler(context.Background(), uc, logger)
	router := handler.InitRouter()
	req, _ := http.NewRequest("GET", "/api/lru", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	var response ResponseItems
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, http.StatusOK, recorder.Code, "status code should be 200")
	assert.Equal(t, expectedResponse.Keys, response.Keys, "Value should match")
}
