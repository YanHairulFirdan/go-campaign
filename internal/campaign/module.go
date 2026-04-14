package campaign

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/app"
	"go-campaign.com/internal/campaign/repository/postgres"
	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/campaign/services"
	v1 "go-campaign.com/internal/campaign/transport/http/v1"
)

func BootHttpV1(router fiber.Router, deps *app.Dependencies) {
	q := sqlc.New(deps.DB) // Now using the shared repository package
	donationRepository := postgres.NewDonationRepository(deps.DB, q)
	campaignRepository := postgres.NewCampaignRepository(q)

	userService := services.NewUserCampaignService(q)
	userHandler := v1.NewHandler(userService)

	publicHandler := v1.NewPublicHandler(
		services.NewCampaignService(
			deps.PaymentGateway,
			donationRepository,
			campaignRepository,
		),
		deps.Config,
	)

	v1.RegisterRoute(router, userHandler, publicHandler)
}
