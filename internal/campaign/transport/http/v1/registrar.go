package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign/repository/postgres"
	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/campaign/services"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/shared/services/payment"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {

	q := sqlc.New(db) // Now using the shared repository package
	donationRepository := postgres.NewDonationRepository(db, q)
	campaignRepository := postgres.NewCampaignRepository(q)

	userService := services.NewUserCampaignService(q)
	userHandler := NewHandler(userService)
	publicHandler := NewPublicHandler(services.NewCampaignService(payment.New(), donationRepository, campaignRepository))

	routeGroup := router.Group("/user/campaigns", middleware.Protected(), middleware.ExtractToken)

	routeGroup.Get("/", userHandler.Index)
	routeGroup.Post("/", userHandler.Create)
	routeGroup.Get("/:id", userHandler.Show)
	routeGroup.Put("/:id", userHandler.Update)

	publicCampaign := router.Group("/campaigns")
	publicCampaign.Get("/", publicHandler.Index)
	publicCampaign.Get("/:slug", publicHandler.Show)
	publicCampaign.Post("/:slug/donate", middleware.Protected(), middleware.ExtractToken, publicHandler.Donate)
	publicCampaign.Get("/:slug/donaturs", publicHandler.Donatur)

	publicCampaign.Post("/xendit/callback", publicHandler.XenditWebhookCallback)
}
