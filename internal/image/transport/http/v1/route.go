package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
)

func RegisterRoute(router fiber.Router, handler *ImageHandler) {
	uploadGroup := router.Group("/images", middleware.Protected(), middleware.ExtractToken)
	uploadGroup.Post("/upload", handler.Upload)
	uploadGroup.Delete("/delete", handler.Delete)
}
