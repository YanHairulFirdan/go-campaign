package campaign

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/campaign/handlers/v1"
	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/campaign/services"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/services/payment"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {

	q := sqlc.New(db) // Now using the shared repository package
	txStore := repository.NewTransactionStore(db)

	userService := services.NewUserCampaignService(q)
	userHandler := v1.NewHandler(userService)
	publicHandler := v1.NewPublicHandler(services.NewCampaignService(q, txStore, payment.New()))

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
	// routeGroup.Delete("/:id", userHandler.Delete)
}
