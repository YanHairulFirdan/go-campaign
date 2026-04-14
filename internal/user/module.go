package user

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/app"
	"go-campaign.com/internal/user/repository/sqlc"
	"go-campaign.com/internal/user/services"
	v1 "go-campaign.com/internal/user/transport/http/v1"
)

func Boot(router fiber.Router, deps *app.Dependencies) {
	q := sqlc.New(deps.DB)
	s := services.NewUserService(q)
	handler := v1.NewHandler(s)

	v1.RegisterRoute(router, handler)
}
