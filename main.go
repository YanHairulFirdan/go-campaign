package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-campaign.com/internal/infrastucture"
	"go-campaign.com/internal/shared/http/middleware"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.RateLimiter())

	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	db, err := infrastucture.InitDatabaseConnection()
	// queries := sqlc.New(db)

	if err != nil {
		panic("Error connecting to the database")
	}

	// txStore := repository.NewTransactionStore(db)

	defer func() {
		if err := db.Close(); err != nil {
			panic("Error closing database connection")
		}
	}()

	infrastucture.InitValidation(db)
	// infrastucture.RegisterRoute(app, queries, txStore, db)
	infrastucture.RegisterRoute(app, db)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = ":3030" // Default port if not set
	}

	err = app.Listen(port)

	if err != nil {
		panic(err)
	}
}
