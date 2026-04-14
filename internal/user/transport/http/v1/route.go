package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
)

func RegisterRoute(router fiber.Router, handler *handler) {
	authV1 := router.Group("/auth")
	authV1.Post("/register", handler.Register)
	authV1.Post("/login", handler.Login)
	authV1.Post("/logout", middleware.Protected(), handler.Logout)
}
