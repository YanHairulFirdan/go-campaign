package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/campaign/services/repository"
	"go-campaign.com/internal/shared/http/request"
)

type CampaignRepository struct {
	sqlc *sqlc.Queries
}

var _ repository.CampaignRepository = (*CampaignRepository)(nil)

func NewCampaignRepository(sqlc *sqlc.Queries) *CampaignRepository {
	return &CampaignRepository{
		sqlc: sqlc,
	}
}

// func (r *CampaignRepository) GetCampaignBySlug(ctx context.Context)
func (r *CampaignRepository) GetPaginatedCampaigns(ctx context.Context, req request.PaginationRequest) ([]repository.CampaignList, error) {
	campaigns, err := r.sqlc.GetCampaigns(ctx, sqlc.GetCampaignsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve campaign list: %w", err)
	}

	var campaignList []repository.CampaignList

	for _, c := range campaigns {
		campaignList = append(campaignList, repository.CampaignList(c))
	}

	return campaignList, nil
}

func (r *CampaignRepository) GetTotalCampaign(ctx context.Context) (int64, error) {
	total, err := r.sqlc.GetTotalCampaigns(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to get total campaign: %w", err)
	}

	return total, nil
}

func (r *CampaignRepository) GetCampaignBySlug(ctx context.Context, slug string) (*repository.DetailCampaign, error) {
	c, err := r.sqlc.GetCampaignBySlug(ctx, slug)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("campaign not found")
		}

		return nil, fmt.Errorf("failed to retrieve the campaign: %w", err)
	}

	return &repository.DetailCampaign{
		ID:            c.ID,
		UserID:        c.UserID,
		Title:         c.Title,
		Description:   c.Description,
		Slug:          c.Slug,
		TargetAmount:  c.TargetAmount,
		CurrentAmount: c.CurrentAmount,
		StartDate:     c.StartDate,
		EndDate:       c.EndDate,
		UserName:      c.UserName,
		UserEmail:     c.UserEmail,
		Progress:      c.Progress,
		Status:        c.Status,
	}, nil
}
