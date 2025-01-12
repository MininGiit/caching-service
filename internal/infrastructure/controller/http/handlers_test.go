
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)
 
type mockUseCase struct {
	result string
}

f

func TestHandlerGet(t *testing.T) {
	ctx := context.Background()
	// 1. Setup mock use case
	mockUseCase := &MockUseCase{
		GetFunc: func(ctx context.Context, key string) (interface{}, time.Time, error) {
			if key == "testKey" {
				return testItem.Value, testItem.ExpiresAt, nil // Return test item
			}
			return nil, time.Time{}, ErrTest // Return error for other keys
		},
	}
	handler := NewHandler(mockUseCase) // Create handler with mock
	router := handler.InitRouter()        // Initialize router

	// 2. Test case 1: Valid key
	t.Run("Valid key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/lru/testKey", nil) // Create a GET request with the key
		recorder := httptest.NewRecorder() // Create a response recorder
		router.ServeHTTP(recorder, req)      // Execute the request

		assert.Equal(t, http.StatusOK, recorder.Code, "Status code should be 200")

		var response ResponseItem // Decode response body into ResponseItem struct
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, "testKey", response.Key, "Key should match")
		assert.Equal(t, testItem.Value, response.Value, "Value should match")
		assert.Equal(t, testItem.ExpiresAt.Format(time.RFC3339), response.ExpiresAt, "Expiration time should match")
	})

	// 3. Test case 2: Non-existent key
	t.Run("Non-existent key", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/lru/nonExistentKey", nil) // Create a GET request with non-existent key
		recorder := httptest.NewRecorder()  // Create a response recorder
		router.ServeHTTP(recorder, req)         // Execute the request

		assert.Equal(t, http.StatusNotFound, recorder.Code, "Status code should be 404")
	})
}


// func TestHandler_get(t *testing.T) {
// 	ctx := context.Background()
// 	// 1. Setup mock use case
// 	mockUseCase := &MockUseCase{
// 		GetFunc: func(ctx context.Context, key string) (interface{}, time.Time, error) {
// 			if key == "testKey" {
// 				return testItem.Value, testItem.ExpiresAt, nil // Return test item
// 			}
// 			return nil, time.Time{}, ErrTest // Return error for other keys
// 		},
// 	}
// 	handler := NewHandler(mockUseCase) // Create handler with mock
// 	router := handler.InitRouter()        // Initialize router

// 	// 2. Test case 1: Valid key
// 	t.Run("Valid key", func(t *testing.T) {
// 		req, _ := http.NewRequest("GET", "/api/lru/testKey", nil) // Create a GET request with the key
// 		recorder := httptest.NewRecorder() // Create a response recorder
// 		router.ServeHTTP(recorder, req)      // Execute the request

// 		assert.Equal(t, http.StatusOK, recorder.Code, "Status code should be 200")

// 		var response ResponseItem // Decode response body into ResponseItem struct
// 		err := json.NewDecoder(recorder.Body).Decode(&response)
// 		assert.NoError(t, err, "Should not return an error")
// 		assert.Equal(t, "testKey", response.Key, "Key should match")
// 		assert.Equal(t, testItem.Value, response.Value, "Value should match")
// 		assert.Equal(t, testItem.ExpiresAt.Format(time.RFC3339), response.ExpiresAt, "Expiration time should match")
// 	})

// 	// 3. Test case 2: Non-existent key
// 	t.Run("Non-existent key", func(t *testing.T) {
// 		req, _ := http.NewRequest("GET", "/api/lru/nonExistentKey", nil) // Create a GET request with non-existent key
// 		recorder := httptest.NewRecorder()  // Create a response recorder
// 		router.ServeHTTP(recorder, req)         // Execute the request

// 		assert.Equal(t, http.StatusNotFound, recorder.Code, "Status code should be 404")
// 	})
// }

// func TestHandler_getAll(t *testing.T) {
// 	ctx := context.Background()
//     // 1. Setup mock use case
// 	mockUseCase := &MockUseCase{
// 		GetAllFunc: func(ctx context.Context) ([]string, []interface{}, error) {
// 			keys := make([]string, len(testItems)) // Create slice of keys
// 			values := make([]interface{}, len(testItems)) // Create slice of values
// 			for i, item := range testItems {
// 				keys[i] = item.Key           // Populate keys
// 				values[i] = item.Value      // Populate values
// 			}
// 			return keys, values, nil // Return values and nil error
// 		},
// 	}
// 	handler := NewHandler(mockUseCase) // Create handler with mock
// 	router := handler.InitRouter()        // Initialize router

// 	// 2. Test case 1: Valid call
// 	t.Run("Valid call", func(t *testing.T) {
// 		req, _ := http.NewRequest("GET", "/api/lru", nil) // Create GET request
// 		recorder := httptest.NewRecorder()           // Create a response recorder
// 		router.ServeHTTP(recorder, req)                  // Execute request

// 		assert.Equal(t, http.StatusOK, recorder.Code, "Status code should be 200")
// 		var response ResponseItems   // Decode response into ResponseItems struct
// 		err := json.NewDecoder(recorder.Body).Decode(&response)
// 		assert.NoError(t, err, "Should not return an error")

// 		keys := make([]string, len(testItems))    // Create expected keys
// 		values := make([]interface{}, len(testItems))    // Create expected values

// 		for i, item := range testItems {
// 			keys[i] = item.Key
// 			values[i] = item.Value
// 		}

// 		assert.Equal(t, keys, response.Keys, "Keys should match") // Validate keys
// 		assert.Equal(t, values, response.Values, "Values should match") // Validate values
// 	})

// 	// 3. Test case 2: Error on GetAll
// 	t.Run("Error on GetAll", func(t *testing.T) {
// 		mockUseCase.GetAllFunc = func(ctx context.Context) ([]string, []interface{}, error) {
// 			return nil, nil, ErrTest // return error
// 		}
// 		req, _ := http.NewRequest("GET", "/api/lru", nil) // Create GET request
// 		recorder := httptest.NewRecorder()           // Create response recorder
// 		router.ServeHTTP(recorder, req)                  // Execute request

// 		assert.Equal(t, http.StatusNotFound, recorder.Code, "Status code should be 404")
// 	})
// }

// func TestHandler_post(t *testing.T) {
//     ctx := context.Background()
//     // 1. Setup mock use case
// 	mockUseCase := &MockUseCase{
// 		PutFunc: func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
// 			if key == "validKey" {
// 				return nil // Return no error for valid key
// 			}
// 			return ErrTest // Return error for invalid key
// 		},
// 	}
// 	handler := NewHandler(mockUseCase) // Create handler
// 	router := handler.InitRouter()        // Initialize router

// 	// 2. Test case 1: Valid request
// 	t.Run("Valid request", func(t *testing.T) {
// 		requestBody := []byte(`{"key": "validKey", "value": "test", "ttlSeconds": 1}`) // Correct json body
// 		req, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody)) // Create post request
// 		recorder := httptest.NewRecorder()
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, http.StatusCreated, recorder.Code, "Status code should be 201")
// 	})

// 	// 3. Test case 2: Invalid request body
// 	t.Run("Invalid request body", func(t *testing.T) {
// 		requestBody := []byte(`invalid json`) // Incorrect JSON
// 		req, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody)) // Create post request
// 		recorder := httptest.NewRecorder()
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Status code should be 400")
// 	})

// 	// 4. Test case 3: Error on Put
// 	t.Run("Error on Put", func(t *testing.T) {
// 		mockUseCase.PutFunc = func(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
// 			return ErrTest // Return error for Put
// 		}
// 		requestBody := []byte(`{"key": "invalidKey", "value": "test", "ttlSeconds": 1}`)
// 		req, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody))
// 		recorder := httptest.NewRecorder()
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Status code should be 400")
// 	})
//     // 5. Test case 4: invalid ttl
// 	t.Run("Invalid ttl", func(t *testing.T) {
// 		requestBody := []byte(`{"key": "validKey", "value": "test", "ttlSeconds": "test"}`) // Incorrect TTL
// 		req, _ := http.NewRequest("POST", "/api/lru", bytes.NewReader(requestBody))
// 		recorder := httptest.NewRecorder()
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Status code should be 400")
// 	})

// }

// func TestHandler_delete(t *testing.T) {
//     ctx := context.Background()
// 	// 1. Setup mock use case
// 	mockUseCase := &MockUseCase{
// 		EvictFunc: func(ctx context.Context, key string) (interface{}, error) {
// 			if key == "testKey" {
// 				return "testValue", nil // return test value for valid key
// 			}
// 			return nil, ErrTest // return error for invalid key
// 		},
// 	}
// 	handler := NewHandler(mockUseCase) // Create handler
// 	router := handler.InitRouter()        // Init router

// 	// 2. Test case 1: Valid delete
// 	t.Run("Valid delete", func(t *testing.T) {
// 		req, _ := http.NewRequest("DELETE", "/api/lru/testKey", nil) // create delete request with testKey
// 		recorder := httptest.NewRecorder()                           // create response recorder
// 		router.ServeHTTP(recorder, req)                            // execute request

// 		assert.Equal(t, http.StatusNoContent, recorder.Code, "Status code should be 204")
// 	})

// 	// 3. Test case 2: Non-existent key delete
// 	t.Run("Non-existent key delete", func(t *testing.T) {
// 		req, _ := http.NewRequest("DELETE", "/api/lru/nonExistentKey", nil) // create delete request with non-existent key
// 		recorder := httptest.NewRecorder()                           // create response recorder
// 		router.ServeHTTP(recorder, req)                            // execute request

// 		assert.Equal(t, http.StatusNotFound, recorder.Code, "Status code should be 404")
// 	})
// }

// func TestHandler_deleteAll(t *testing.T) {
// 	ctx := context.Background()
//     // 1. Setup mock use case
// 	mockUseCase := &MockUseCase{
// 		EvictAllFunc: func(ctx context.Context) error {
// 			return nil // Mock use case success
// 		},
// 	}

// 	handler := NewHandler(mockUseCase) // Create handler
// 	router := handler.InitRouter()        // Initialize router

// 	// 2. Test case 1: Valid deleteAll
// 	t.Run("Valid deleteAll", func(t *testing.T) {
// 		req, _ := http.NewRequest("DELETE", "/api/lru", nil) // Create Delete request to /api/lru
// 		recorder := httptest.NewRecorder()
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, http.StatusNoContent, recorder.Code, "Status code should be 204")
// 	})

//     // 3. Test case 2: Error on EvictAll
//     t.Run("Error on EvictAll", func(t *testing.T) {
//         mockUseCase.EvictAllFunc = func(ctx context.Context) error {
//             return ErrTest // Mock use case fail
//         }
//         req, _ := http.NewRequest("DELETE", "/api/lru", nil)
//         recorder := httptest.NewRecorder()
//         router.ServeHTTP(recorder, req)
//         assert.Equal(t, http.StatusNoContent, recorder.Code, "Status code should be 204")
//     })
// }
