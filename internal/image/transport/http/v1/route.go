package v1

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRouteV1(apiV1 fiber.Router) {
	uploadGroup := apiV1.Group("/images")
	uploadGroup.Post("/upload", Upload)
	uploadGroup.Delete("/delete", Delete)
}
