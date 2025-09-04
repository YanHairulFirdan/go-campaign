package infrastucture

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	campaignV1 "go-campaign.com/internal/campaign/transport/http/v1"
	"go-campaign.com/internal/image"
	userV1 "go-campaign.com/internal/user/transport/http/v1"
)

func RegisterRoute(app *fiber.App, db *sql.DB) {
	apiV1 := app.Group("/api/v1")

	userV1.RegisterRouteV1(apiV1, db)
	campaignV1.RegisterRouteV1(apiV1, db)
	image.RegisterRouteV1(apiV1)
}
