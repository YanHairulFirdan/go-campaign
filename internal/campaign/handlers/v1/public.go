package v1

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/internal/shared/repository/sqlc"
)

type publicHandler struct {
	q *sqlc.Queries
}

func NewPublicHandler(q *sqlc.Queries) *publicHandler {
	return &publicHandler{
		q: q,
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
