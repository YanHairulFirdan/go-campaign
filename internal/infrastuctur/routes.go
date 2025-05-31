package infrastuctur

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/internal/user"
)

func RegisterRoute(app *fiber.App, db *sql.DB, q *sqlc.Queries) {
	apiV1 := app.Group("/api/v1")

	user.RegisterRouteV1(apiV1, q)
	campaign.RegisterRouteV1(apiV1, db)
}
