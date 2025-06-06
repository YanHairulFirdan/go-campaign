package v1

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign/entities"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	q *sqlc.Queries
}

// NewHandler creates a new user handler with the given repository.
func NewHandler(q *sqlc.Queries) *handler {
	return &handler{
		q: q,
	}
}

// listCampaigns lists all campaigns for a user.
func (h *handler) Index(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	pb := response.NewPaginationBuilder(
		perPage,
		page,
		func() ([]sqlc.GetPaginatedUserCampaignRow, error) {
			campaigns, err := h.q.GetPaginatedUserCampaign(c.Context(), sqlc.GetPaginatedUserCampaignParams{
				UserID: int32(userID),
				Limit:  int32(perPage),
				Offset: int32((page - 1) * perPage),
				Title:  c.Query("title", ""),
				Status: int32(c.QueryInt("status", int(entities.StatusActive))),
			})

			if err != nil {
				return nil, err
			}

			return campaigns, nil
		},
		func() (int, error) {
			totalCount, err := h.q.GetTotalUserCampaigns(c.Context(), sqlc.GetTotalUserCampaignsParams{
				UserID: int32(userID),
				Title:  c.Query("title", ""),
				Status: int32(c.QueryInt("status", int(entities.StatusActive))),
			})
			if err != nil {
				return 0, err
			}
			return int(totalCount), nil
		},
	)

	pagination, err := pb.Build()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(pagination)
}

// createCampaign creates a new campaign for a user.
func (h *handler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

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

	campaign, err := h.q.CreateCampaign(c.Context(), sqlc.CreateCampaignParams{
		UserID:       int32(userID),
		Title:        req.Title,
		Description:  &req.Description,
		Slug:         req.Slug,
		TargetAmount: strconv.FormatFloat(float64(req.TargetAmount), 'f', -1, 32),
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       int32(req.Status),
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
	userID := c.Locals("userID").(int)
	campaignID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid campaign ID",
				"Campaign ID must be a valid integer",
			),
		)
	}

	campaign, err := h.q.GetUserCampaignById(c.Context(), sqlc.GetUserCampaignByIdParams{
		ID:     int32(campaignID),
		UserID: int32(userID),
	})

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(
		response.NewResponse(
			"success",
			"Campaign retrieved successfully",
			campaign,
		),
	)
}

func (h *handler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	campaignID, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid campaign ID",
				"Campaign ID must be a valid integer",
			),
		)
	}

	cp, err := h.q.GetUserCampaignById(c.Context(), sqlc.GetUserCampaignByIdParams{
		ID:     int32(campaignID),
		UserID: int32(userID),
	})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
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
	cp.Description = &req.Description
	cp.Slug = req.Slug
	cp.TargetAmount = strconv.FormatFloat(float64(req.TargetAmount), 'f', -1, 32)
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

	cp.Status = int32(entities.Status(req.Status))

	updatedCampaign, err := h.q.UpdateCampaign(c.Context(), sqlc.UpdateCampaignParams{
		Title:        cp.Title,
		Description:  cp.Description,
		Slug:         cp.Slug,
		TargetAmount: cp.TargetAmount,
		StartDate:    cp.StartDate,
		EndDate:      cp.EndDate,
		Status:       cp.Status,
		UserID:       cp.UserID,
		ID:           cp.ID,
	})
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
			updatedCampaign,
		),
	)
}

// func (h *handler) Delete(c *fiber.Ctx) error {
// 	userID, err := auth.ValidateToken(c.Locals("user").(*jwt.Token).Raw)

// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Unauthorized",
// 				"Failed to validate token",
// 			),
// 		)
// 	}

// 	campaignID := c.Params("id")

// 	cp, err := h.r.FindBy("id", campaignID)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Campaign not found",
// 				err.Error(),
// 			),
// 		)
// 	}

// 	if cp.UserID != userID {
// 		return c.Status(fiber.StatusForbidden).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Forbidden",
// 				"You do not have permission to delete this campaign",
// 			),
// 		)
// 	}

// 	deletedAt := time.Now()
// 	cp.DeletedAt = &deletedAt

// 	_, err = h.r.Update(cp)

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Failed to delete campaign",
// 				err.Error(),
// 			),
// 		)
// 	}

// 	return c.Status(200).JSON(
// 		response.NewResponse(
// 			"success",
// 			"Campaign deleted successfully",
// 			nil,
// 		),
// 	)
// }
