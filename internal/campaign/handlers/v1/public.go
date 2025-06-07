package v1

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/internal/shared/services/payment"
	"go-campaign.com/pkg/auth"
	"go-campaign.com/pkg/validation"
)

type publicHandler struct {
	q       *sqlc.Queries
	txStore *repository.TransactionStore
	pg      payment.PaymentGateway
}

func NewPublicHandler(
	q *sqlc.Queries,
	txStore *repository.TransactionStore,
	pg payment.PaymentGateway,
) *publicHandler {
	return &publicHandler{
		q:       q,
		txStore: txStore,
		pg:      pg,
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

	pb := response.NewPaginationBuilder(
		perPage,
		page,
		func() ([]sqlc.GetCampaignsRow, error) {
			campaigns, err := h.q.GetCampaigns(c.Context(), sqlc.GetCampaignsParams{
				Limit:  int32(perPage),
				Offset: int32((page - 1) * perPage),
			})

			if err != nil {
				return nil, err
			}

			return campaigns, nil
		},
		func() (int, error) {
			count, err := h.q.GetTotalCampaigns(c.Context())
			if err != nil {
				return 0, err
			}

			return int(count), nil
		},
	)

	campaigns, err := pb.Build()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Internal server error",
				err.Error(),
			),
		)
	}

	return c.Status(200).JSON(campaigns)
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

	campaign, err := h.q.GetCampaignBySlug(c.Context(), slug)
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
	jwtToken := c.Locals("user").(*jwt.Token)

	userID, ok := auth.ValidateToken(jwtToken.Raw)

	var url string

	if ok != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse(
				"error",
				"Unauthorized",
				"Invalid token",
			),
		)
	}

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

	campaign, err := h.q.FindCampaignsBySlugForUpdate(c.Context(), slug)
	if err != nil || campaign.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			response.NewErrorResponse(
				"error",
				"Campaign not found",
				"Campaign with the provided slug does not exist",
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

	err = h.txStore.ExecTx(func() error {
		// create donaturs
		d, err := h.q.CreateDonatur(c.Context(), sqlc.CreateDonaturParams{
			CampaignID: campaign.ID,
			UserID:     int32(userID),
			Name:       donationRequest.Name,
			Email: sql.NullString{
				String: donationRequest.Email,
				Valid:  donationRequest.Email != "",
			},
		})

		if err != nil {
			return fmt.Errorf("failed to create donatur: %w", err)
		}

		donation, err := h.q.CreateDonation(c.Context(), sqlc.CreateDonationParams{
			DonaturID:  d.ID,
			CampaignID: campaign.ID,
			Amount:     fmt.Sprintf("%d", donationRequest.Amount),
			Note: sql.NullString{
				String: donationRequest.Note,
				Valid:  donationRequest.Note != "",
			},
		})

		if err != nil {
			return fmt.Errorf("failed to create donation: %w", err)
		}

		transactionID := uuid.New()
		invoiceRequest := payment.InvoiceRequest{
			ExternalID: transactionID,
			Amount:     float64(donationRequest.Amount),
			Currency:   "IDR",
			UserDetail: payment.UserDetail{
				Email:    donationRequest.Email,
				FullName: donationRequest.Name,
			},
			ProductDetails: []payment.ProductDetail{
				{Name: "Donation to campaig", Price: float64(donationRequest.Amount), Quantity: 1},
			},
		}

		url, err = h.pg.CreateInvoice(invoiceRequest)

		if err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}

		// create payment
		_, err = h.q.CreatePayment(c.Context(), sqlc.CreatePaymentParams{
			TransactionID: transactionID,
			DonaturID:     d.ID,
			DonationID:    donation.ID,
			CampaignID:    campaign.ID,
			Amount:        fmt.Sprintf("%d", donationRequest.Amount),
			Link: sql.NullString{
				String: url,
			},
			Note: sql.NullString{
				String: donationRequest.Note,
			},
			Status: int32(sqlc.DonationPaymentStatusPending),
		})

		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		return nil
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

	// return c.Status(fiber.StatusOK).JSON(
	// 	response.NewResponse(
	// 		"success",
	// 		"Webhook callback received successfully",
	// 		nil,
	// 	),
	// )

	h.txStore.ExecTx(func() error {
		// Find the payment by transaction ID
		p, err := h.q.GetPaymentByTransactionId(c.Context(), uuid.MustParse(webhookEvent.ExternalID))

		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("payment not found for transaction ID: %s", webhookEvent.ExternalID)
			}

			return fmt.Errorf("failed to get payment: %w", err)
		}

		jsonResponse, err := webhookEvent.ToJson()

		if err != nil {

			return fmt.Errorf("failed to convert webhook event to JSON: %w", err)
		}

		if webhookEvent.Status != payment.InvoiceStatusPaid {
			_, err = h.q.UpdatePaymentFromCallback(c.Context(), sqlc.UpdatePaymentFromCallbackParams{
				Status: int32(sqlc.DonationPaymentStatusExpired),
				ID:     p.ID,
				Vendor: sql.NullString{
					String: "xendit",
				},
				Method: sql.NullString{
					String: webhookEvent.PaymentMethod,
				},
				Response: pqtype.NullRawMessage{
					RawMessage: []byte(jsonResponse),
					Valid:      true,
				},
				PaymentDate: sql.NullTime{
					Time: func() time.Time {
						parsedTime, err := time.Parse("2006-01-02T15:04:05Z", webhookEvent.PaidAt)
						if err != nil {
							return time.Time{}
						}
						return parsedTime
					}(),
					Valid: webhookEvent.PaidAt != "",
				},
			})

			if err != nil {
				return fmt.Errorf("failed to update payment status to expired: %w", err)
			}

			return nil
		}

		_, err = h.q.UpdatePaymentFromCallback(c.Context(), sqlc.UpdatePaymentFromCallbackParams{
			Status: int32(sqlc.DonationPaymentStatusPaid),
			ID:     p.ID,
			Vendor: sql.NullString{
				String: "xendit",
			},
			Method: sql.NullString{
				String: webhookEvent.PaymentMethod,
			},
			Response: pqtype.NullRawMessage{
				RawMessage: []byte(jsonResponse),
				Valid:      true,
			},
			PaymentDate: sql.NullTime{
				Time: func() time.Time {
					parsedTime, err := time.Parse("2006-01-02T15:04:05Z", webhookEvent.PaidAt)
					if err != nil {
						return time.Time{}
					}
					return parsedTime
				}(),
				Valid: webhookEvent.PaidAt != "",
			},
		})

		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		log.Printf("Status: %s, Amount: %.2f, Paid Amount: %.2f, Campaign ID: %d",
			webhookEvent.Status,
			webhookEvent.Amount,
			webhookEvent.PaidAmount,
			p.CampaignID,
		)

		cp, err := h.q.FindCampaignByIdForUpdate(c.Context(), p.CampaignID)

		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("campaign not found for ID: %d", p.CampaignID)
			}

			return fmt.Errorf("failed to find campaign: %w", err)
		}

		err = h.q.Donate(c.Context(), sqlc.DonateParams{
			ID:     cp.ID,
			Amount: fmt.Sprintf("%.2f", webhookEvent.Amount),
		})

		if err != nil {
			return fmt.Errorf("failed to update campaign amount: %w", err)
		}

		return nil
	})

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

	pb := response.NewPaginationBuilder(
		perPage,
		page,
		func() ([]sqlc.GetPaginatedDonatursRow, error) {
			donaturs, err := h.q.GetPaginatedDonaturs(c.Context(), sqlc.GetPaginatedDonatursParams{
				Slug:   slug,
				Limit:  int32(perPage),
				Offset: int32((page - 1) * perPage),
			})
			if err != nil {
				return nil, err
			}

			return donaturs, nil
		},
		func() (int, error) {
			count, err := h.q.GetCampaignTotalPaidDonaturs(c.Context(), slug)

			if err != nil {
				return 0, err
			}

			return int(count), nil
		},
	)

	donaturs, err := pb.Build()

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
			"Donaturs retrieved successfully",
			donaturs,
		),
	)
}
