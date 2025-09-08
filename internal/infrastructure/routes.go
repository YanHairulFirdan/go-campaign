package infrastructure

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	campaignV1 "go-campaign.com/internal/campaign/transport/http/v1"
	imageV1 "go-campaign.com/internal/image/transport/http/v1"
	userV1 "go-campaign.com/internal/user/transport/http/v1"
)

func RegisterRoute(app *fiber.App, db *sql.DB) {
	apiV1 := app.Group("/api/v1")

	userV1.RegisterRouteV1(apiV1, db)
	campaignV1.RegisterRouteV1(apiV1, db)
	imageV1.RegisterRouteV1(apiV1)
}
