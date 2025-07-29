package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/middleware"
	v1 "go-campaign.com/internal/user/handlers/v1"
	"go-campaign.com/internal/user/repository/sqlc"
	"go-campaign.com/internal/user/services"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {
	q := sqlc.New(db)
	s := services.NewUserService(q)
	handler := v1.NewHandler(s)

	authV1 := router.Group("/auth")
	authV1.Post("/register", handler.Register)
	authV1.Post("/login", handler.Login)
	authV1.Post("/logout", middleware.Protected(), handler.Logout)
}
