package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/user/handlers/v1"
	"go-campaign.com/internal/user/repository"
)

func RegisterRouteV1(router fiber.Router, db *sql.DB) {
	handler := v1.NewHandler(repository.NewRepository(db))

	authV1 := router.Group("/auth")
	authV1.Post("/register", handler.Register)
	authV1.Post("/login", handler.Login)

}
