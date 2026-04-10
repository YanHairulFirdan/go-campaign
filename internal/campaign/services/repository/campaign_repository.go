package repository

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"go-campaign.com/internal/shared/http/request"
)

type CampaignRepository interface {
	GetPaginatedCampaigns(ctx context.Context, req request.PaginationRequest) ([]CampaignList, error)
	GetTotalCampaign(ctx context.Context) (int64, error)
	GetCampaignBySlug(ctx context.Context, slug string) (*DetailCampaign, error)
}

type CampaignList struct {
	ID            int32           `json:"id"`
	Title         string          `json:"title"`
	Slug          string          `json:"slug"`
	CurrentAmount decimal.Decimal `json:"current_amount"`
	TargetAmount  decimal.Decimal `json:"target_amount"`
	Progress      decimal.Decimal `json:"progress"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       time.Time       `json:"end_date"`
	Status        string          `json:"status"`
}

type DetailCampaign struct {
	ID            int32           `json:"id"`
	Title         string          `json:"title"`
	Description   *string         `json:"description"`
	Slug          string          `json:"slug"`
	TargetAmount  decimal.Decimal `json:"target_amount"`
	CurrentAmount decimal.Decimal `json:"current_amount"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       time.Time       `json:"end_date"`
	UserName      string          `json:"user_name"`
	UserEmail     string          `json:"user_email"`
	Progress      decimal.Decimal `json:"progress"`
	Status        int32           `json:"status"`
}
