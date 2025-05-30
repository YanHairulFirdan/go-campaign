package campaign

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/campaign/handlers/v1"
	"go-campaign.com/internal/campaign/repository"
	"go-campaign.com/internal/shared/http/middleware"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {
	userHandler := v1.NewHandler(repository.NewRepository(db))

	routeGroup := router.Group("/user/campaigns", middleware.Protected())

	routeGroup.Get("/", userHandler.Index)
	routeGroup.Post("/", userHandler.Create)
	routeGroup.Get("/:id", userHandler.Show)
	routeGroup.Put("/:id", userHandler.Update)
	routeGroup.Delete("/:id", userHandler.Delete)
}
