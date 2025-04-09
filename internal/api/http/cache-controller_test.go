package http

import (
	"bytes"
	"cache_engine_httpserver/internal/api/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

// CreateCache is the function you're testing
func CreateCacheData(c *fiber.Ctx, ctx *model.CacheAppContext) error {
	// Your implementation here
	return nil
}

func TestCreateCache(t *testing.T) {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route that uses the CreateCache function
	app.Post("http://127.0.0.1:3000/cache-engine-api/create", func(c fiber.Ctx) error {
		// Initialize your custom CacheAppContext
		cacheCtx := &model.CacheAppContext{
			// Populate fields as necessary
		}
		return CreateCacheData(&c, cacheCtx)
	})

	// Define test cases
	tests := []struct {
		name           string
		requestBody    []byte
		expectedStatus int
	}{
		{
			name:           "Valid Request",
			requestBody:    []byte(`{"key":"example","value":"data", "duration_in_seconds":10}`), // Adjust based on your expected input
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Request",
			requestBody:    []byte(`{"invalid":"data"}`), // Adjust based on your validation logic
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP request
			req := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:3000/cache-engine-api/create", bytes.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform the request using Fiber's app.Test method
			resp, err := app.Test(req, -1) // -1 disables the timeout
			assert.NoError(t, err)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Println("Body :", bodyString)

			// Assert the response status code
			// assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Additional assertions can be added here, such as checking the response body
		})
	}
}
