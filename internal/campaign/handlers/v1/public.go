package v1

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/campaign/services"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/services/payment"
	"go-campaign.com/pkg/validation"
)

type publicHandler struct {
	s *services.CampaignService
}

func NewPublicHandler(
	s *services.CampaignService,
) *publicHandler {
	return &publicHandler{
		s: s,
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

	campaigns, totalCount, err := h.s.GetCampaigns(c.Context(), int32(page), int32(perPage))

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

	campaign, err := h.s.FindCampaignsBySlugForUpdate(c.Context(), slug)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				err.Error(),
			),
		)
	}

	if userID == int(campaign.ID) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Cannot donate to your own campaign",
				"You cannot donate to your own campaign",
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

	errMessages, err := validation.Validate(donationRequest, nil)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Validation error",
				err.Error(),
			),
		)
	}

	if len(errMessages) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewValidationErrorResponse(
				"error",
				"Validation error",
				errMessages,
			),
		)
	}

	url, err := h.s.Donate(c.Context(), services.DonationRequest{
		CampaignID: campaign.ID,
		UserID:     int32(userID),
		Amount:     donationRequest.Amount,
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
	// return c.Status(fiber.StatusOK).JSON(
	// 	response.NewResponse(
	// 		"success",
	// 		"Webhook callback received successfully",
	// 		nil,
	// 	),
	// )
	var webhookEvent payment.XenditInvoiceWebhookResponse

	if err := c.BodyParser(&webhookEvent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid request body",
				"Failed to parse request body",
			),
		)
	}

	if err := h.s.UpdatePaymentFromCallback(c.Context(), webhookEvent); err != nil {
		log.Printf("Error updating payment from callback: %v", err)
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

// func (h *publicHandler) Donatur(c *fiber.Ctx) error {
// 	slug := c.Params("slug")
// 	if slug == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Slug is required",
// 				"Slug parameter cannot be empty",
// 			),
// 		)
// 	}

// 	page := c.QueryInt("page", 1)
// 	if page <= 0 {
// 		page = 1
// 	}

// 	perPage := c.QueryInt("per_page", 10)
// 	if perPage <= 0 {
// 		perPage = 10
// 	}

// 	pb := response.NewPaginationBuilder(
// 		perPage,
// 		page,
// 		func() ([]sqlc.GetPaginatedDonatursRow, error) {
// 			donaturs, err := h.q.GetPaginatedDonaturs(c.Context(), sqlc.GetPaginatedDonatursParams{
// 				Slug:   slug,
// 				Limit:  int32(perPage),
// 				Offset: int32((page - 1) * perPage),
// 			})
// 			if err != nil {
// 				return nil, err
// 			}

// 			return donaturs, nil
// 		},
// 		func() (int, error) {
// 			count, err := h.q.GetCampaignTotalPaidDonaturs(c.Context(), slug)

// 			if err != nil {
// 				return 0, err
// 			}

// 			return int(count), nil
// 		},
// 	)

// 	donaturs, err := pb.Build()

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(
// 			response.NewErrorResponse(
// 				"error",
// 				"Internal server error",
// 				err.Error(),
// 			),
// 		)
// 	}

// 	return c.Status(200).JSON(
// 		response.NewResponse(
// 			"success",
// 			"Donaturs retrieved successfully",
// 			donaturs,
// 		),
// 	)
// }
