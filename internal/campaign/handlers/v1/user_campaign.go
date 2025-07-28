package v1

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign/entities"
	"go-campaign.com/internal/campaign/services"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/pkg/validation"
)

type handler struct {
	s *services.UserCampaignService
}

// NewHandler creates a new user handler with the given repository.
func NewHandler(s *services.UserCampaignService) *handler {
	return &handler{
		s: s,
	}
}

// listCampaigns lists all campaigns for a user.
func (h *handler) Index(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	title := c.Query("title", "")
	status := c.QueryInt("status", int(entities.StatusDraft))
	campaigns, totalCount, err := h.s.GetPaginatedUserCampaigns(
		c.Context(), services.PaginatedCampaignRequest{
			UserID: int32(userID),
			Limit:  int32(perPage),
			Title:  title,
			Status: int32(status),
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

	return c.Status(200).JSON(response.NewPagination(
		"success",
		"Campaigns retrieved successfully",
		campaigns,
		response.NewMeta(
			page,
			perPage,
			int(totalCount),
		)),
	)
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

	uploaded, err := c.MultipartForm()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid multipart form data",
				err.Error(),
			),
		)
	}

	if uploaded == nil || len(uploaded.File) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Missing image file",
				"Image file is required",
			),
		)
	}

	req.Images = uploaded.File["images"]

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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewValidationErrorResponse(
				"error",
				"Validation failed",
				validationErrors,
			),
		)
	}

	uploadDir := "./public/uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewErrorResponse(
					"error",
					"Failed to create upload directory",
					err.Error(),
				),
			)
		}
	}

	images := make([]string, 0, len(req.Images))
	// looping through all uploaded images and saving them
	for _, fileHeader := range req.Images {
		extension := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):])
		fileHeader.Filename = fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		images = append(images, fileHeader.Filename)
		c.SaveFile(fileHeader, fmt.Sprintf("./%s/%s", uploadDir, fileHeader.Filename))
	}

	campaign, err := h.s.CreateCampaign(c.Context(), services.CreateCampaignRequest{
		UserID:       int32(userID),
		Title:        req.Title,
		Description:  req.Description,
		Slug:         req.Slug,
		TargetAmount: req.TargetAmount,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Status:       int(req.Status),
		Images:       images,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to create campaign",
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

	campaign, err := h.s.FindUserCampaign(c.Context(), int32(userID), int32(campaignID))

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

	cp, err := h.s.FindUserCampaign(c.Context(), int32(userID), int32(campaignID))
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

	updatedCampaign, err := h.s.UpdateCampaign(c.Context(), cp.ID, services.CreateCampaignRequest{
		Title:        req.Title,
		Description:  req.Description,
		Slug:         req.Slug,
		TargetAmount: req.TargetAmount,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Status:       int(req.Status),
		UserID:       int32(userID),
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
