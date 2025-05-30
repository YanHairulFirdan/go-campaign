package v1

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go-campaign.com/cmd/api/response"
	"go-campaign.com/internal/campaign/entities"
	"go-campaign.com/internal/campaign/repository"
	"go-campaign.com/pkg/auth"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	r repository.Repository
}

// NewHandler creates a new user handler with the given repository.
func NewHandler(r repository.Repository) *handler {
	return &handler{
		r: r,
	}
}

// listCampaigns lists all campaigns for a user.
func (h *handler) Index(c *fiber.Ctx) error {
	jwtToken := c.Locals("user").(*jwt.Token)

	userID, ok := auth.ValidateToken(jwtToken.Raw)

	if ok != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Invalid user ID in token",
			),
		)
	}

	page, err := strconv.Atoi(c.Query("page", "1"))

	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid page number",
				"Page number must be a positive integer",
			),
		)
	}
	perPage, err := strconv.Atoi(c.Query("per_page", "10"))
	if err != nil || perPage < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid per_page number",
				"Per page number must be a positive integer",
			),
		)
	}

	campaigns, err := h.r.Paginate(repository.Filters{
		{
			Column:   "user_id",
			Value:    userID,
			Operator: "=",
		},
	}, page, perPage)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"Campaigns retrieved successfully",
			listCampaignCollection(campaigns),
		),
	) // Placeholder return
}

// createCampaign creates a new campaign for a user.
func (h *handler) Create(c *fiber.Ctx) error {
	jwtToken := c.Locals("user").(*jwt.Token)

	userID, ok := auth.ValidateToken(jwtToken.Raw)

	if ok != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Invalid user ID in token",
			),
		)
	}

	var req createCampaignRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid request body",
				err.Error(),
			),
		)
	}

	validationErrors, err := validation.Validate(req, nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewValidationErrorResponse(
				"error",
				"Validation failed",
				validationErrors,
			),
		)
	}

	startDate, err := time.Parse(time.DateTime, req.StartDate)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid start date format",
				"Start date must be in date-time format (e.g., 2023-10-01 12:00:00)",
			),
		)
	}

	endDate, err := time.Parse(time.DateTime, req.EndDate)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid end date format",
				"End date must be in date-time format (e.g., 2023-10-01 12:00:00)",
			),
		)
	}

	campaign, err := h.r.Create(entities.Campaign{
		UserID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		Slug:          req.Slug,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0, // Initial current amount is 0
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        entities.Status(req.Status),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	return c.Status(201).JSON(
		response.NewResponse(
			"success",
			"Campaign created successfully",
			map[string]interface{}{
				"id":             campaign.ID,
				"title":          campaign.Title,
				"description":    campaign.Description,
				"slug":           campaign.Slug,
				"target_amount":  campaign.TargetAmount,
				"current_amount": campaign.CurrentAmount,
				"start_date":     campaign.StartDate.Format(time.RFC3339),
				"end_date":       campaign.EndDate.Format(time.RFC3339),
				"status":         campaign.Status,
			},
		),
	)
}

func (h *handler) Show(c *fiber.Ctx) error {
	jwtToken := c.Locals("user").(*jwt.Token)

	userID, err := auth.ValidateToken(jwtToken.Raw)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Failed validate token",
			),
		)
	}

	campaignID := c.Params("id")
	campaign, err := h.r.FindBy("id", campaignID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}

	if campaign.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			response.NewErrorResponse(
				"error",
				"Forbidden",
				"You do not have permission to access this campaign",
			),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"Campaign retrieved successfully",
			map[string]any{
				"id":             campaign.ID,
				"title":          campaign.Title,
				"description":    campaign.Description,
				"slug":           campaign.Slug,
				"target_amount":  campaign.TargetAmount,
				"current_amount": campaign.CurrentAmount,
				"start_date":     campaign.StartDate.Format(time.RFC3339),
				"end_date":       campaign.EndDate.Format(time.RFC3339),
				"status":         campaign.Status,
			},
		),
	)
}

func (h *handler) Update(c *fiber.Ctx) error {
	userID, err := auth.ValidateToken(c.Locals("user").(*jwt.Token).Raw)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Failed to validate token",
			),
		)
	}

	campaignID := c.Params("id")
	cp, err := h.r.FindBy("id", campaignID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}
	if cp.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			response.NewErrorResponse(
				"error",
				"Forbidden",
				"You do not have permission to update this campaign",
			),
		)
	}

	var req updateCampaignRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid request body",
				err.Error(),
			),
		)
	}

	validationErrors, err := validation.Validate(req, validation.ValidationExceptions{
		"Slug": {
			Column: "id",
			Value:  campaignID,
		},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewValidationErrorResponse(
				"error",
				"Validation failed",
				validationErrors,
			),
		)
	}

	// Update the campaign with the new data
	cp.Title = req.Title
	cp.Description = req.Description
	cp.Slug = req.Slug
	cp.TargetAmount = req.TargetAmount
	cp.StartDate, err = time.Parse(time.DateTime, req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid start date format",
				"Start date must be in date-time format (e.g., 2023-10-01 12:00:00)",
			),
		)
	}
	cp.EndDate, err = time.Parse(time.DateTime, req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid end date format",
				"End date must be in date-time format (e.g., 2023-10-01 12:00:00)",
			),
		)
	}

	cp.Status = entities.Status(req.Status)

	updatedCampaign, err := h.r.Update(cp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to update campaign",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"Campaign updated successfully",
			map[string]any{
				"id":             updatedCampaign.ID,
				"title":          updatedCampaign.Title,
				"description":    updatedCampaign.Description,
				"slug":           updatedCampaign.Slug,
				"target_amount":  updatedCampaign.TargetAmount,
				"current_amount": updatedCampaign.CurrentAmount,
				"start_date":     updatedCampaign.StartDate.Format(time.RFC3339),
				"end_date":       updatedCampaign.EndDate.Format(time.RFC3339),
				"status":         updatedCampaign.Status,
			},
		),
	)
}

func (h *handler) Delete(c *fiber.Ctx) error {
	userID, err := auth.ValidateToken(c.Locals("user").(*jwt.Token).Raw)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Failed to validate token",
			),
		)
	}

	campaignID := c.Params("id")

	cp, err := h.r.FindBy("id", campaignID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}

	if cp.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(
			response.NewErrorResponse(
				"error",
				"Forbidden",
				"You do not have permission to delete this campaign",
			),
		)
	}

	deletedAt := time.Now()
	cp.DeletedAt = &deletedAt

	_, err = h.r.Update(cp)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to delete campaign",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"Campaign deleted successfully",
			nil,
		),
	)
}
