package v1

import (
	"fmt"
	"log"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go-campaign.com/internal/campaign/entities"
	"go-campaign.com/internal/campaign/services"
	"go-campaign.com/internal/config"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/services/payment"
	"go-campaign.com/pkg/validation"
)

type publicHandler struct {
	s      *services.CampaignService
	config *config.Config
}

func NewPublicHandler(
	s *services.CampaignService,
	c *config.Config,

) *publicHandler {
	return &publicHandler{
		s:      s,
		config: c,
	}
}

func (h *publicHandler) Index(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)

	if page <= 0 {
		page = 1
	}

	if perPage <= 0 {
		perPage = 10
	}

	campaigns, totalCount, err := h.s.GetCampaigns(
		c.Context(),
		(int32(page)-1)*int32(perPage),
		int32(perPage),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"status":  "error",
				"message": "Internal server error",
				"error":   err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewPagination(
			"success",
			"Campaigns retrieved successfully",
			campaigns,
			response.NewMeta(
				page,
				perPage,
				totalCount,
			),
		),
	)
}

// show by slug
func (h *publicHandler) Show(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Slug is required",
				"Slug parameter cannot be empty",
			),
		)
	}

	log.Println("Fetching campaign with slug:", slug)

	campaign, err := h.s.GetCampaignBySlug(c.Context(), slug)
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

func (h *publicHandler) Donate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Slug is required",
				"Slug parameter cannot be empty",
			),
		)
	}

	campaign, err := h.s.GetCampaignBySlug(c.Context(), slug)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}

	if userID == int(campaign.UserID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Cannot donate to your own campaign",
				"You cannot donate to your own campaign",
			),
		)
	}

	if campaign.Status != int32(entities.StatusActive) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign is not active",
				"Cannot donate to a non-active campaign",
			),
		)
	}

	var donationRequest DonationRequest

	if err := c.BodyParser(&donationRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid request body",
				"Failed to parse request body",
			),
		)
	}

	err = donationRequest.Validate()
	validationErr, err := validation.ParseValidationErrors(err)
	if err != nil {
		return c.Status(500).JSON(
			response.NewErrorResponse("error", "Internal server error", err.Error()),
		)
	}

	if len(validationErr) > 0 {
		return c.Status(422).JSON(
			response.NewFailedValidationErrorResponse("error", "Validation failed", validationErr),
		)
	}

	url, err := h.s.Donate(c.Context(), services.DonationRequest{
		CampaignID: campaign.ID,
		UserID:     int32(userID),
		Amount:     decimal.NewFromFloat32(donationRequest.Amount),
		Name:       donationRequest.Name,
		Email:      donationRequest.Email,
		Note:       &donationRequest.Note,
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

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(
			"success",
			"Donation successful",
			fiber.Map{
				"message": "Thank you for your donation! You will be redirected to the payment page shortly.",
				"link":    url,
			},
		),
	)
}

func (h *publicHandler) XenditWebhookCallback(c *fiber.Ctx) error {
	webhookEvent, err := parseWebhookResponse(c, h.config.App.Service.Payment.Vendor)

	if err != nil {
		return fmt.Errorf("failed parsing response: %w", err)
	}

	_, exists := payment.MapPaymentStatus[webhookEvent.Status]

	if !exists {
		return fmt.Errorf("payment status is invalid: %w", webhookEvent.Status)
	}

	if err := h.s.UpdatePaymentFromCallback(c.Context(), webhookEvent); err != nil {
		log.Printf("Error updating payment from callback: %w", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(
			"success",
			"Webhook callback received successfully",
			nil,
		),
	)
}

func (h *publicHandler) Donatur(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Slug is required",
				"Slug parameter cannot be empty",
			),
		)
	}

	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	perPage := c.QueryInt("per_page", 10)
	if perPage <= 0 {
		perPage = 10
	}

	donaturs, totalCount, err := h.s.GetDonatur(c.Context(), services.GetDonaturListRequest{
		Slug:   slug,
		Limit:  int32(perPage),
		Offset: int32((page - 1) * perPage),
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

	return c.Status(200).JSON(
		response.NewPagination(
			"success",
			"Donaturs retrieved successfully",
			donaturs,
			response.NewMeta(
				page,
				perPage,
				totalCount,
			),
		),
	)
}

func parseWebhookResponse(c *fiber.Ctx, vendor string) (*payment.PaymentCallback, error) {
	availablePaymentVendor := []string{"xendit"}

	if !slices.Contains(availablePaymentVendor, vendor) {
		return nil, fmt.Errorf("currently we do not support payment from requested vendor: %s", vendor)
	}
	var webhookEvent payment.XenditInvoiceWebhookResponse

	if err := c.BodyParser(&webhookEvent); err != nil {
		return nil, fmt.Errorf("failed to parse webhook body: %w", err)
	}

	rawData, err := webhookEvent.ToJson()

	if err != nil {
		return nil, fmt.Errorf("failed to convert webhook event to JSON: %w", err)
	}

	return &payment.PaymentCallback{
		ExternalID:    webhookEvent.ExternalID,
		RawData:       rawData,
		PaidAt:        webhookEvent.PaidAt,
		Status:        payment.PaymentStatus(webhookEvent.Status),
		PaymentMethod: webhookEvent.PaymentMethod,
	}, nil
}
