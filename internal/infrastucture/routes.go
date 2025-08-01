package infrastucture

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign"
	"go-campaign.com/internal/image"
	"go-campaign.com/internal/user"
)

func RegisterRoute(app *fiber.App, db *sql.DB) {
	apiV1 := app.Group("/api/v1")

	user.RegisterRouteV1(apiV1, db)
	campaign.RegisterRouteV1(apiV1, db)
	image.RegisterRouteV1(apiV1)
}
