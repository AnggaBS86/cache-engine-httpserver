package http

import (
	"cache_engine_httpserver/internal/api/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

// Set up a simple Fiber app
func setUpAppCache(key string, value string) *fiber.App {
	app := fiber.New()
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	ctx := new(model.CacheAppContext)
	ctx.Cache = cache

	expiration := time.Now().Add(time.Duration(111) * time.Second)
	entry := model.CacheEntry{
		Value:      value,
		Expiration: expiration,
	}
	entryData, _ := json.Marshal(entry)
	_ = ctx.Cache.Set(key, entryData)

	// Define a route
	app.Get("cache-engine-api/get", func(c fiber.Ctx) error {
		key := c.Query("key")
		return c.JSON(fiber.Map{
			"status": "OK",
			"cache": fiber.Map{
				"key":   key,
				"value": entry.Value,
			},
		})
	})

	return app
}

type Cache struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Cache  Cache  `json:"cache"`
	Status string `json:"status"`
}

func TestGetHandler(t *testing.T) {
	cacheKey := "key"
	cacheValue := "1234"

	app := setUpAppCache(cacheKey, cacheValue)
	req := httptest.NewRequest("GET", "/cache-engine-api/get?key="+cacheKey, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error occurred while making request: %v", err)
	}

	// Assert the status code is 200 OK
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	var responseData Response
	err = json.Unmarshal(bodyBytes, &responseData)

	if cacheKey != responseData.Cache.Key {
		t.Error("Test error")
	}

	assert.Equal(t, cacheKey, responseData.Cache.Key)
	assert.Equal(t, cacheValue, responseData.Cache.Value)
	assert.Equal(t, "OK", responseData.Status)
}

func TestGetHandlerWithEmptyCacheKey(t *testing.T) {
	cacheKey := ""
	cacheValue := "123456"

	app := setUpAppCache(cacheKey, cacheValue)
	req := httptest.NewRequest("GET", "/cache-engine-api/get?"+cacheKey+"="+cacheKey, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error occurred while making request: %v", err)
	}

	// Assert the status code is 200 OK
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	var responseData Response
	err = json.Unmarshal(bodyBytes, &responseData)

	if cacheKey != responseData.Cache.Key {
		t.Error("Test error")
	}

	assert.Equal(t, cacheKey, "")
	assert.Equal(t, cacheValue, responseData.Cache.Value)
	assert.Equal(t, "OK", responseData.Status)
}

func TestGetHandlerWithEmptyCacheValue(t *testing.T) {
	cacheKey := "key"
	cacheValue := ""

	app := setUpAppCache(cacheKey, cacheValue)
	req := httptest.NewRequest("GET", "/cache-engine-api/get?"+cacheKey+"="+cacheKey, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error occurred while making request: %v", err)
	}

	// Assert the status code is 200 OK
	assert.Equal(t, 200, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	var responseData Response
	err = json.Unmarshal(bodyBytes, &responseData)

	if cacheKey != responseData.Cache.Key {
		t.Error("Test error")
	}

	assert.Equal(t, cacheKey, responseData.Cache.Key)
	assert.Equal(t, cacheValue, "")
	assert.Equal(t, "OK", responseData.Status)
}

// Define the payload structure
type Payload struct {
	Key               string `json:"key"`
	Value             string `json:"value"`
	DurationInSeconds int    `json:"duration_in_seconds"`
}

// Define the response structure
type CacheResponse struct {
	Cache struct {
		Key               string `json:"key"`
		Value             string `json:"value"`
		DurationInSeconds int    `json:"duration_in_seconds"`
	} `json:"cache"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Handler to test
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create the response
	response := CacheResponse{
		Message: "Value set successfully",
		Status:  "OK",
	}
	response.Cache.Key = payload.Key
	response.Cache.Value = payload.Value
	response.Cache.DurationInSeconds = payload.DurationInSeconds

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func TestCreatingCache(t *testing.T) {
	payload := `{
		"key": "username",
		"value": "Angga",
		"duration_in_seconds": 10
	}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler
	postHandler(recorder, req)

	// Verify the response
	resp := recorder.Result()
	defer resp.Body.Close()

	// Assert status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Decode the response body into the CacheResponse struct
	var response CacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	// Validate "cache" object
	if response.Cache.Key != "username" {
		t.Errorf("expected key %q, got %q", "username", response.Cache.Key)
	}
	if response.Cache.Value != "Angga" {
		t.Errorf("expected value %q, got %q", "Angga", response.Cache.Value)
	}
	if response.Cache.DurationInSeconds != 10 {
		t.Errorf("expected duration %d, got %d", 10, response.Cache.DurationInSeconds)
	}

	// Validate "message"
	if response.Message != "Value set successfully" {
		t.Errorf("expected message %q, got %q", "Value set successfully", response.Message)
	}

	// Validate "status"
	if response.Status != "OK" {
		t.Errorf("expected status %q, got %q", "OK", response.Status)
	}
}

func TestCreatingCacheWithCacheValueIsEmpty(t *testing.T) {
	payload := `{
		"key": "username",
		"value": "",
		"duration_in_seconds": 10
	}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler
	postHandler(recorder, req)

	// Verify the response
	resp := recorder.Result()
	defer resp.Body.Close()

	// Assert status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Decode the response body into the CacheResponse struct
	var response CacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	// Validate "cache" object
	if response.Cache.Key != "username" {
		t.Errorf("expected key %q, got %q", "username", response.Cache.Key)
	}
	if response.Cache.Value != "" {
		t.Errorf("expected value %q, got %q", "", response.Cache.Value)
	}
	if response.Cache.DurationInSeconds != 10 {
		t.Errorf("expected duration %d, got %d", 10, response.Cache.DurationInSeconds)
	}

	// Validate "message"
	if response.Message != "Value set successfully" {
		t.Errorf("expected message %q, got %q", "Value set successfully", response.Message)
	}

	// Validate "status"
	if response.Status != "OK" {
		t.Errorf("expected status %q, got %q", "OK", response.Status)
	}
}

func TestCreatingCacheWithInvalidDurationParam(t *testing.T) {
	payload := `{
		"key": "username",
		"value": "",
		"duration_in_seconds": -1
	}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	// Create a recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler
	postHandler(recorder, req)

	// Verify the response
	resp := recorder.Result()
	defer resp.Body.Close()

	// Assert status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Decode the response body into the CacheResponse struct
	var response CacheResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	fmt.Println("cache response : ", response)

	// Validate "cache" object
	if response.Cache.Key != "username" {
		t.Errorf("expected key %q, got %q", "username", response.Cache.Key)
	}
	if response.Cache.Value != "" {
		t.Errorf("expected value %q, got %q", "", response.Cache.Value)
	}
	if response.Cache.DurationInSeconds != -1 {
		t.Errorf("expected duration %d, got %d", -1, response.Cache.DurationInSeconds)
	}

	// Validate "message"
	if response.Message != "Value set successfully" {
		t.Errorf("expected message %q, got %q", "Value set successfully", response.Message)
	}

	// Validate "status"
	if response.Status != "OK" {
		t.Errorf("expected status %q, got %q", "OK", response.Status)
	}
}
