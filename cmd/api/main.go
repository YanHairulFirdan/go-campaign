package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go-campaign.com/cmd/api/v1/auth"
)

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apiV1 := app.Group("/api/v1")

	authV1 := apiV1.Group("/auth")

	authV1handler := auth.NewHandler()
	authV1.Post("/register", authV1handler.Register)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = ":3030" // Default port if not set
	}

	err = app.Listen(port)

	if err != nil {
		panic(err)
	}
}
