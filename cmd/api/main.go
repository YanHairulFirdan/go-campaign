package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-campaign.com/cmd/api/v1/auth"
	"go-campaign.com/internal/user/repository"
	"go-campaign.com/pkg/validation"
)

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}

	dbConnectionString := os.Getenv("DATABASE_CONNECTION")

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))
	}
	err = validation.Init(db)

	defer db.Close()

	if err != nil {
		panic(fmt.Sprintf("Error initializing validation: %v", err))
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apiV1 := app.Group("/api/v1")

	authV1 := apiV1.Group("/auth")

	userRepo := repository.NewRepository(db)
	authV1handler := auth.NewHandler(userRepo)

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
