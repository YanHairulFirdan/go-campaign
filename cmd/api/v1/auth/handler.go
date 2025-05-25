package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/cmd/api/response"
	"go-campaign.com/internal/user"
	"go-campaign.com/internal/user/repository"
	"go-campaign.com/pkg/hash"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	repository repository.Repository
}

func NewHandler(repository repository.Repository) *handler {
	return &handler{
		repository: repository,
	}
}

func (h *handler) Register(c *fiber.Ctx) error {
	var req UserRegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(
			response.NewErrorResponse("error", "Invalid request body", err.Error()),
		)
	}

	validationErrors, err := validation.Validate(req)

	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Internal server error", err.Error()),
		)
	}

	if len(validationErrors) > 0 {
		return c.Status(422).JSON(
			response.NewValidationErrorResponse("error", "Validation failed", validationErrors),
		)
	}

	password, err := hash.Password(req.Password)
	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Failed to hash password", err.Error()),
		)
	}

	user := user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	_, err = h.repository.Create(user)

	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Failed to create user", err.Error()),
		)
	}

	return c.Status(200).JSON(response.NewResponse("success", "User registered successfully", nil))
}
