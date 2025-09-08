package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"go-campaign.com/internal/config"
	"go-campaign.com/internal/infrastructure"
	"go-campaign.com/internal/shared/http/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	deps, err := setupDependencies(cfg)

	if err != nil {
		return fmt.Errorf("error setting up dependencies: %v", err)
	}

	port := cfg.App.Port
	log.Printf("Starting server on port %s in %s mode...", port, cfg.App.ENV)

	defer deps.Close()

	app := setupApp(deps)

	go func() {
		if err := app.Listen(port); err != nil {
			log.Printf("error starting server: %v", err)
		}
	}()

	<-c
	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	return nil
}

type Dependencies struct {
	DB     *sql.DB
	Config *config.Config
}

func (d *Dependencies) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}

	return nil
}

func setupDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := infrastructure.InitDatabaseConnection(cfg)

	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	infrastructure.InitValidation(db)

	return &Dependencies{
		DB:     db,
		Config: cfg,
	}, nil
}

func setupApp(deps *Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})
	app.Use(logger.New())
	app.Use(middleware.RateLimiter())
	app.Static("/", "./public")

	infrastructure.RegisterRoute(app, deps.DB)

	return app
}
