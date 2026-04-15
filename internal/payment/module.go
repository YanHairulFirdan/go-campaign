package payment

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/payment/repository/postgres"
	"go-campaign.com/internal/payment/repository/sqlc"
	"go-campaign.com/internal/payment/service"
	v1 "go-campaign.com/internal/payment/transport/http/v1"
)

type HTTPDeps struct {
	DB *sql.DB
}

func BootHttpV1(fiberApp fiber.Router, deps HTTPDeps) {
	q := sqlc.New(deps.DB)
	r := postgres.NewPaymentRepository(q)
	s := service.NewPaymentService(r)
	h := v1.NewPaymentHandler(s)

	v1.RegisterRoute(fiberApp, h)
}
