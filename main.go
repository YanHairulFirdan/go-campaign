package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-campaign.com/internal/infrastuctur"
	"go-campaign.com/internal/shared/repository/sqlc"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	db, err := infrastuctur.InitDatabaseConnection()
	queries := sqlc.New(db)

	if err != nil {
		panic("Error connecting to the database")
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic("Error closing database connection")
		}
	}()

	infrastuctur.InitValidation(db)
	infrastuctur.RegisterRoute(app, queries)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = ":3030" // Default port if not set
	}

	err = app.Listen(port)

	if err != nil {
		panic(err)
	}
}
