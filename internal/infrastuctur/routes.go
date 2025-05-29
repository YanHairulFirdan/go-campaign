package infrastuctur

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/cmd/api/middleware"
	"go-campaign.com/cmd/api/v1/auth"
	"go-campaign.com/cmd/api/v1/user"
	campaignRepo "go-campaign.com/internal/campaign/repository"
	userRepo "go-campaign.com/internal/user/repository"
)

func RegisterRoute(app *fiber.App, db *sql.DB) {
	ur := userRepo.NewRepository(db)
	cr := campaignRepo.NewRepository(db)

	apiV1 := app.Group("/api/v1")

	authV1 := apiV1.Group("/auth")
	authV1handler := auth.NewHandler(ur)
	authV1.Post("/register", authV1handler.Register)
	authV1.Post("/login", authV1handler.Login)

	campaignV1Handler := user.NewHandler(cr)
	campaignV1 := apiV1.Group("/user/campaigns", middleware.Protected())
	campaignV1.Get("/", campaignV1Handler.Index)
	campaignV1.Post("/", campaignV1Handler.Create)
	campaignV1.Get("/:id", campaignV1Handler.Show)
	campaignV1.Put("/:id", campaignV1Handler.Update)
	campaignV1.Delete("/:id", campaignV1Handler.Delete)
}
