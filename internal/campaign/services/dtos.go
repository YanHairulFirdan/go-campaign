package services

import "time"

type PaginatedCampaignRequest struct {
	UserID int32
	Limit  int32
	Title  string
	Status int32
}

type CreateCampaignRequest struct {
	UserID       int32
	Title        string
	Description  string
	Slug         string
	TargetAmount float32
	StartDate    string
	EndDate      string
	Status       int
}

type Campaign struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Slug          string    `json:"slug"`
	TargetAmount  float32   `json:"target_amount"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        int       `json:"status"` // 0: Draft, 1: Active
	CurrentAmount float32   `json:"current_amount"`
}

type DonationRequest struct {
	CampaignID int32
	UserID     int32
	Amount     float32
	Name       string
	Email      string
	Note       *string
}
