package repository

import (
	"database/sql"

	"go-campaign.com/internal/campaign"
)

type Filter struct {
	Column   string
	Value    any
	Operator string
}

type Filters []Filter

type Repository interface {
	Create(campaign.Campaign) (campaign.Campaign, error)
	FindBy(column string, value any) (campaign.Campaign, error)
	GetCampaignsFromUser(userID int) ([]campaign.Campaign, error)
	Update(campaign.Campaign) (campaign.Campaign, error)
}

func NewRepository(connection *sql.DB) Repository {
	return newPostgresRepository(connection)
}
