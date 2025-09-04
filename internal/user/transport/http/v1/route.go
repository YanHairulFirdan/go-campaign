package v1

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
	"go-campaign.com/internal/user/repository/sqlc"
	"go-campaign.com/internal/user/services"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {
	q := sqlc.New(db)
	s := services.NewUserService(q)
	handler := NewHandler(s)

	authV1 := router.Group("/auth")
	authV1.Post("/register", handler.Register)
	authV1.Post("/login", handler.Login)
	authV1.Post("/logout", middleware.Protected(), handler.Logout)
}
