package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               20,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use the client's IP address as the key for rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "You have exceeded the rate limit. Please try again later.",
			})
		},
	})
}
