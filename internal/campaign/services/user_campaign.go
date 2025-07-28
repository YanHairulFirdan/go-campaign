package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-campaign.com/internal/campaign/repository/sqlc"
)

type UserCampaignService struct {
	q *sqlc.Queries
}

func NewUserCampaignService(q *sqlc.Queries) *UserCampaignService {
	return &UserCampaignService{
		q: q,
	}
}

func (s *UserCampaignService) GetPaginatedUserCampaigns(ctx context.Context, request PaginatedCampaignRequest) ([]sqlc.GetPaginatedUserCampaignRow, int64, error) {
	campaigns, err := s.q.GetPaginatedUserCampaign(
		ctx,
		sqlc.GetPaginatedUserCampaignParams{
			UserID: request.UserID,
			Limit:  request.Limit,
			Offset: (request.UserID - 1) * request.Limit,
			Title:  request.Title,
			Status: request.Status,
		},
	)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated user campaigns: %w", err)
	}

	totalCount, err := s.q.GetTotalUserCampaigns(ctx, sqlc.GetTotalUserCampaignsParams{
		UserID: request.UserID,
		Title:  request.Title,
		Status: request.Status,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total user campaigns: %w", err)
	}

	return campaigns, totalCount, nil
}

func (s *UserCampaignService) CreateCampaign(ctx context.Context, request CreateCampaignRequest) (*sqlc.Campaign, error) {
	startDate, err := time.Parse(time.DateTime, request.StartDate)

	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse(time.DateTime, request.EndDate)

	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	log.Print(request.Description)

	campaign, err := s.q.CreateCampaign(ctx, sqlc.CreateCampaignParams{
		UserID:       request.UserID,
		Title:        request.Title,
		Description:  &request.Description,
		Slug:         request.Slug,
		TargetAmount: &request.TargetAmount,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       int32(request.Status),
		Images:       request.Images,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return &campaign, nil
}

func (s *UserCampaignService) FindUserCampaign(ctx context.Context, userID, campaignID int32) (*sqlc.GetUserCampaignByIdRow, error) {
	c, err := s.q.GetUserCampaignById(ctx, sqlc.GetUserCampaignByIdParams{
		ID:     campaignID,
		UserID: userID,
	})

	if err != nil {
		return nil, fmt.Errorf("Campaign not found: %w", err)
	}

	return &c, nil
}

func (s *UserCampaignService) UpdateCampaign(ctx context.Context, campaignID int32, request CreateCampaignRequest) (*sqlc.UpdateCampaignRow, error) {
	startDate, err := time.Parse(time.DateTime, request.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse(time.DateTime, request.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	campaign, err := s.q.UpdateCampaign(ctx, sqlc.UpdateCampaignParams{
		ID:           campaignID,
		UserID:       request.UserID,
		Title:        request.Title,
		Description:  &request.Description,
		Slug:         request.Slug,
		TargetAmount: &request.TargetAmount,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       int32(request.Status),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return &campaign, nil
}
