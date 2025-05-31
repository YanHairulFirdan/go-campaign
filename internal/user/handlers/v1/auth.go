package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/internal/user/entities"
	"go-campaign.com/pkg/auth"
	"go-campaign.com/pkg/hash"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	queries *sqlc.Queries
}

func NewHandler(q *sqlc.Queries) *handler {
	return &handler{
		queries: q,
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

	password, err := hash.Password(req.Password)
	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Failed to hash password", err.Error()),
		)
	}

	user := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	_, err = h.queries.CreateUser(c.Context(), sqlc.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Failed to create user", err.Error()),
		)
	}

	jwtToken, err := auth.GenerateToken(user.ID)

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

	user, err := h.queries.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.Status(404).JSON(
			response.NewErrorResponse("error", "User not found", err.Error()),
		)
	}

	if match, err := hash.ComparePassword(user.Password, req.Password); err != nil || !match {
		return c.Status(401).JSON(
			response.NewErrorResponse("error", "Invalid credentials", "Email & Password does not match"),
		)
	}

	jwtToken, err := auth.GenerateToken(int(user.ID))
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
