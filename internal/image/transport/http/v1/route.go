package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
)

func RegisterRouteV1(apiV1 fiber.Router) {
	uploadGroup := apiV1.Group("/images", middleware.Protected(), middleware.ExtractToken)
	uploadGroup.Post("/upload", Upload)
	uploadGroup.Delete("/delete", Delete)
}
