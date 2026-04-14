package infrastructure

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/config"
)

func RegisterRoute(app *fiber.App, db *sql.DB, config config.Config) error {
	// apiV1 := app.Group("/api/v1")

	// userV1.RegisterRouteV1(apiV1, db)
	// err := campaignV1.RegisterRouteV1(apiV1, db, config)

	// if err != nil {
	// 	return err
	// }

	// imageV1.RegisterRouteV1(apiV1)

	return nil
}
