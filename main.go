package main

import (
	"cache_engine_httpserver/internal/api/middleware"
	"cache_engine_httpserver/internal/api/model"
	"cache_engine_httpserver/internal/api/router"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func getDefaultCacheDuration() time.Duration {
	defaultCacheDurationInSeconds, err := strconv.Atoi(os.Getenv("DEFAULT_CACHE_DURATION_IN_SECONDS"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	return time.Duration(defaultCacheDurationInSeconds) * time.Second
}

// Define a separate function for the middleware
func firstHandler(c fiber.Ctx) error {
	fmt.Println("🥇 First handler")
	return c.Next()
}

func main() {
	// Load env configuration from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err.Error())
	}

	// Initiliaze cache
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(getDefaultCacheDuration()))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create AppContext to share dependencies
	appContext := &model.CacheAppContext{
		Cache: cache,
	}

	// Initialize Fiber app
	app := fiber.New()
	app.Use(firstHandler)
	// Or extend your config for customization
	app.Use(middleware.RateLimiterMiddleware())
	router.HandleRoute(app, appContext)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(port))
}
