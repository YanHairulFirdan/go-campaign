package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
)

func RegisterRoute(router fiber.Router, userHandler *handler, publicHandler *publicHandler) error {
	routeGroup := router.Group("/user/campaigns", middleware.Protected(), middleware.ExtractToken)

	routeGroup.Get(
		"/",
		middleware.PaginationQueryNormalizer(middleware.QueryNormalization{
			"page":     1,
			"per_page": 10,
		}),
		userHandler.Index,
	)
	routeGroup.Post("/", userHandler.Create)
	routeGroup.Get("/:id", userHandler.Show)
	routeGroup.Put("/:id", userHandler.Update)

	publicCampaign := router.Group("/campaigns")
	publicCampaign.Get(
		"/",
		middleware.PaginationQueryNormalizer(middleware.QueryNormalization{
			"page":     1,
			"per_page": 10,
		}),
		publicHandler.Index,
	)
	publicCampaign.Get("/:slug", publicHandler.Show)
	publicCampaign.Post("/:slug/donate", middleware.Protected(), middleware.ExtractToken, publicHandler.Donate)
	publicCampaign.Get("/:slug/donaturs", publicHandler.Donatur)

	publicCampaign.Post("/xendit/callback", publicHandler.XenditWebhookCallback)

	return nil
}
