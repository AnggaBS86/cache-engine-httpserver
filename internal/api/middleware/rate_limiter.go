package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// Function to return the limiter middleware
func RateLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return c.IP() == "127.0.0.1" // Allow localhost to bypass rate limiting
		},
		Max:        1020,             // Max requests allowed
		Expiration: 30 * time.Second, // Expiration time window
		KeyGenerator: func(c fiber.Ctx) string {
			return c.Get("x-forwarded-for") // Use "X-Forwarded-For" header
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.SendFile("./toofast.html") // Custom response when limit is reached
		},
	})
}
