package image

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/image/handlers"
)

func RegisterRouteV1(apiV1 fiber.Router) {
	uploadGroup := apiV1.Group("/images")
	uploadGroup.Post("/upload", handlers.Upload)
	uploadGroup.Delete("/delete", handlers.Delete)
}
