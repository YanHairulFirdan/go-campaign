package auth

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/cmd/api/response"
	"go-campaign.com/pkg/validation"
)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
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

	return c.Status(200).JSON(response.NewResponse("success", "User registered successfully", nil))
}
