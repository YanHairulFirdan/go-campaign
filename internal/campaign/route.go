package campaign

import (
	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/campaign/handlers/v1"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/repository/sqlc"
)

func RegisterRouteV1(router fiber.Router, q *sqlc.Queries, txStore *repository.TransactionStore) {
	userHandler := v1.NewHandler(q)
	publicHandler := v1.NewPublicHandler(q, txStore)

	routeGroup := router.Group("/user/campaigns", middleware.Protected(), middleware.ExtractToken)

	routeGroup.Get("/", userHandler.Index)
	routeGroup.Post("/", userHandler.Create)
	routeGroup.Get("/:id", userHandler.Show)
	routeGroup.Put("/:id", userHandler.Update)

	publicCampaign := router.Group("/campaigns")
	publicCampaign.Get("/", publicHandler.Index)
	publicCampaign.Get("/:slug", publicHandler.Show)
	publicCampaign.Post("/:slug/donate", middleware.Protected(), publicHandler.Donate)
	// routeGroup.Delete("/:id", userHandler.Delete)
}
