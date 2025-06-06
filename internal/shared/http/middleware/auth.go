package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go-campaign.com/pkg/auth"
)

func Protected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
	})
}

func ExtractToken(c *fiber.Ctx) error {
	jwtToken := c.Locals("user").(*jwt.Token)

	if jwtToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "No token provided",
		})
	}

	userID, ok := auth.ValidateToken(jwtToken.Raw)
	if ok != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Invalid token",
		})
	}
	c.Locals("userID", userID)
	return c.Next()
}
