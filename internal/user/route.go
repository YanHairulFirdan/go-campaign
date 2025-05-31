package user

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/repository/sqlc"
	v1 "go-campaign.com/internal/user/handlers/v1"
)

func RegisterRouteV1(router fiber.Router, q *sqlc.Queries) {
	handler := v1.NewHandler(q)

	authV1 := router.Group("/auth")
	authV1.Post("/register", handler.Register)
	authV1.Post("/login", handler.Login)

}
