package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-campaign.com/cmd/api/middleware"
	"go-campaign.com/cmd/api/v1/auth"
	"go-campaign.com/internal/user/repository"
	"go-campaign.com/pkg/validation"
)

type App struct {
	DB *sql.DB
}

func newApp() *App {
	dbConnectionString := os.Getenv("DATABASE_CONNECTION")

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))
	}
	err = validation.Init(db)

	if err != nil {
		panic(fmt.Sprintf("Error initializing validation: %v", err))
	}

	return &App{
		DB: db,
	}
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	App := newApp()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	apiV1 := app.Group("/api/v1")

	authV1 := apiV1.Group("/auth")

	userRepo := repository.NewRepository(App.DB)
	authV1handler := auth.NewHandler(userRepo)

	apiV1.Get("/protected", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"user":    c.Locals("user"),
		})
	})

	// open
	apiV1.Get("/open", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "This is an open route",
			"status":  "success",
		})
	})

	authV1.Post("/register", authV1handler.Register)
	authV1.Post("/login", authV1handler.Login)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = ":3030" // Default port if not set
	}

	err = app.Listen(port)

	defer App.DB.Close()

	if err != nil {
		panic(err)
	}
}
