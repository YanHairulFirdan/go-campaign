package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/user/repository/sqlc"
	"go-campaign.com/internal/user/services"
	v1 "go-campaign.com/internal/user/transport/http/v1"
)

type HTTPDeps struct {
	DB *sql.DB
}

func BootHttpV1(router fiber.Router, deps HTTPDeps) {
	q := sqlc.New(deps.DB)
	s := services.NewUserService(q)
	handler := v1.NewHandler(s)

	v1.RegisterRoute(router, handler)
}
