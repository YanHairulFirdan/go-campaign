package infrastuctur

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/internal/user"
)

func RegisterRoute(app *fiber.App, q *sqlc.Queries, txStore *repository.TransactionStore) {
	apiV1 := app.Group("/api/v1")

	user.RegisterRouteV1(apiV1, q)
	campaign.RegisterRouteV1(apiV1, q, txStore)
}
