package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/user/services"
	"go-campaign.com/pkg/auth"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	s *services.UserService
}

func NewHandler(s *services.UserService) *handler {
	return &handler{
		s: s,
	}
}

func (h *handler) Register(c *fiber.Ctx) error {
	var req UserRegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(
			response.NewErrorResponse("error", "Invalid request body", err.Error()),
		)
	}

	validationErrors, err := validation.Validate(req, nil)

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

	userID, err := h.s.CreateUser(c.Context(), services.CreateUserDTO{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Failed to create user", err.Error()),
		)
	}

	jwtToken, err := auth.GenerateToken(int(userID))

	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Error when generating token", err.Error()),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"User registered successfully",
			map[string]string{
				"token": jwtToken,
			}),
	)
}

func (h *handler) Login(c *fiber.Ctx) error {
	var req UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(
			response.NewErrorResponse("error", "Invalid request body", err.Error()),
		)
	}

	validationErrors, err := validation.Validate(req, nil)

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

	userID, err := h.s.CheckLoginUser(c.Context(), req.Email, req.Password)

	if err != nil || userID == 0 {
		return c.Status(401).JSON(
			response.NewErrorResponse("error", "Invalid email or password", "Unauthorized"),
		)

	}

	jwtToken, err := auth.GenerateToken(userID)
	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Error when generating token", err.Error()),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"User logged in successfully",
			map[string]string{
				"token": jwtToken,
			},
		),
	)
}

func (h *handler) Logout(c *fiber.Ctx) error {
	// Clear the JWT token from the cookie
	c.ClearCookie("jwt_token")

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"User logged out successfully",
			nil,
		),
	)
}
