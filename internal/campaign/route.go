package campaign

import (
	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/campaign/handlers/v1"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/shared/repository/sqlc"
)

func RegisterRouteV1(router fiber.Router, q *sqlc.Queries) {
	userHandler := v1.NewHandler(q)
	publicHandler := v1.NewPublicHandler(q)

	routeGroup := router.Group("/user/campaigns", middleware.Protected())

	routeGroup.Get("/", userHandler.Index)
	routeGroup.Post("/", userHandler.Create)
	routeGroup.Get("/:id", userHandler.Show)
	routeGroup.Put("/:id", userHandler.Update)

	publicCampaign := router.Group("/campaigns")
	publicCampaign.Get("/", publicHandler.Index)
	// routeGroup.Delete("/:id", userHandler.Delete)
}
