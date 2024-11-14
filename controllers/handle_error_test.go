package controllers

import (
	"errors"
	"github.com/spf13/cast"
	"net/http"
	"testing"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		errorCode      int
		err            error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Test CODE_NOT_FOUND",
			errorCode:      CODE_NOT_FOUND,
			err:            errors.New("test error"),
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error_msg":  "not_found",
				"error_code": CODE_NOT_FOUND,
				"detail":     "test error",
			},
		},
		{
			name:           "Test CODE_TIMEOUT",
			errorCode:      CODE_TIMEOUT,
			err:            errors.New("CODE_TIMEOUT"),
			expectedStatus: http.StatusGatewayTimeout,
			expectedBody: map[string]interface{}{
				"error_msg":  "timeout",
				"error_code": CODE_TIMEOUT,
				"detail":     "CODE_TIMEOUT",
			},
		},
		{
			name:           "Test CODE_INVALID_PARAMS",
			errorCode:      CODE_INVALID_PARAMS,
			err:            errors.New("invalid params"),
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error_msg":  "invalid params",
				"error_code": CODE_INVALID_PARAMS,
				"detail":     "invalid params",
			},
		},
	}

	// Setup Gin router for testing
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin context
			r := gin.Default()
			r.POST("/test", func(c *gin.Context) {
				handleError(c, tt.errorCode, tt.err)
			})

			// Perform the request
			w := performRequest(r, "POST", "/test", tt.err)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedBody["error_msg"], responseBody["error_msg"])
			assert.Equal(t, tt.expectedBody["error_code"], cast.ToInt(responseBody["error_code"]))
			assert.Equal(t, tt.expectedBody["detail"], responseBody["detail"])
		})
	}
}

func TestHandleSuccess(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Test handle success with data",
			data:           map[string]string{"message": "transfer successful"},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error_msg":  "OK",
				"error_code": CODE_SUCCESS,
				"data":       map[string]interface{}(map[string]interface{}{"message": "transfer successful"}),
			},
		},
		{
			name:           "Test handle success with nil data",
			data:           nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"error_msg":  "OK",
				"error_code": CODE_SUCCESS,
				"data":       nil,
			},
		},
	}

	// Setup Gin router for testing
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin context
			r := gin.Default()
			r.POST("/test", func(c *gin.Context) {
				handleSuccess(c, tt.data)
			})

			// Perform the request
			w := performRequest(r, "POST", "/test", tt.data)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			var responseBody map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedBody["error_msg"], responseBody["error_msg"])
			assert.Equal(t, tt.expectedBody["error_code"], cast.ToInt(responseBody["error_code"]))
			assert.Equal(t, tt.expectedBody["data"], responseBody["data"])
		})
	}
}

// Helper function to perform HTTP requests in tests
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	// Create a new request
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
