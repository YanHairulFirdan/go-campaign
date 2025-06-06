package v1

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/repository/sqlc"
	"go-campaign.com/pkg/auth"
	"go-campaign.com/pkg/validation"
)

type publicHandler struct {
	q       *sqlc.Queries
	txStore *repository.TransactionStore
}

func NewPublicHandler(
	q *sqlc.Queries,
	txStore *repository.TransactionStore,
) *publicHandler {
	return &publicHandler{
		q:       q,
		txStore: txStore,
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

		_, err = h.q.CreateDonation(c.Context(), sqlc.CreateDonationParams{
			DonaturID:  d.ID,
			CampaignID: campaign.ID,
			Amount:     fmt.Sprintf("%d", donationRequest.Amount),
			Note: sql.NullString{
				String: donationRequest.Note,
				Valid:  donationRequest.Note != "",
			},
			PaymentStatus: int32(sqlc.DonationPaymentStatusPending),
		})

		if err != nil {
			return fmt.Errorf("failed to create donation: %w", err)
		}

		if err != nil {
			return fmt.Errorf("failed to donate: %w", err)
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
			nil,
		),
	)
}
