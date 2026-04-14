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
	"go-campaign.com/internal/app"
	"go-campaign.com/internal/campaign"
	"go-campaign.com/internal/config"
	"go-campaign.com/internal/image"
	"go-campaign.com/internal/infrastructure"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/shared/services/payment"
	"go-campaign.com/internal/user"
	"go-campaign.com/pkg/filesystem"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %w", err)
	}
}

func run() error {
	c := make(chan os.Signal, 1)
	serverErr := make(chan error, 1)
	signal.Notify(c, os.Interrupt)

	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("failed to validate the configuration: %w", err)
	}

	deps, err := setupDependencies(cfg)

	if err != nil {
		return fmt.Errorf("error setting up dependencies: %w", err)
	}

	port := cfg.App.Port
	log.Printf("Starting server on port %s in %s mode...", port, cfg.App.ENV)

	defer deps.CloseDatabaseConnection()

	app, err := setupApp(deps)

	if err != nil {
		return fmt.Errorf("failed to setup the application: %w", err)
	}

	setupModule(app, deps)

	go func() {
		serverErr <- app.Listen(port)
	}()

	select {
	case <-c:
		log.Println("Gracefully shutting down...")
	case err := <-serverErr:
		return fmt.Errorf("start server: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("server shutdown error: %w", err)
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

func setupDependencies(cfg *config.Config) (*app.Dependencies, error) {
	db, err := infrastructure.InitDatabaseConnection(cfg)

	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	infrastructure.InitValidation(db)

	fsystem := filesystem.NewLocalFileSystem()
	paymentGateway, err := payment.New(cfg.App.Service.Payment.SecretKey)

	if err != nil {
		return nil, err
	}

	return &app.Dependencies{
		DB:             db,
		Config:         cfg,
		FileSystem:     fsystem,
		PaymentGateway: paymentGateway,
	}, nil
}

func setupApp(deps *app.Dependencies) (*fiber.App, error) {
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

	return app, nil
}

func setupModule(fiberApp *fiber.App, deps *app.Dependencies) {
	modules := []app.Bootable{
		campaign.BootHttpV1,
		user.BootHttpV1,
		image.BootHttpV1,
	}

	v1 := fiberApp.Group("api/v1")

	for _, module := range modules {
		module(v1, deps)
	}
}
