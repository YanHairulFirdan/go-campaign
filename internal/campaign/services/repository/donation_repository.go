package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DonationRepository interface {
	CreateDonationIntent(ctx context.Context, req CreateDonationIntentParams) (*DonationIntent, error)
	MarkInvoiceCreated(ctx context.Context, paymentID int32, link, vendor string) error
	MarkInvoiceFailed(ctx context.Context, paymentID int32) error
	UpdateDonationPaymentFromWebhook(ctx context.Context, req UpdatePayment) error
	GetPaginatedDonatur(ctx context.Context, req GetPaginatedDonaturParams) ([]DonaturList, error)
	GetTotalPaidDonatur(ctx context.Context, slug string) (int64, error)
}

type CreateDonationIntentParams struct {
	CampaignID int32
	UserID     int32
	Amount     decimal.Decimal
	Name       string
	Email      string
	Note       *string
}

type DonationIntent struct {
	PaymentID     int32
	TransactionID uuid.UUID
	CampaignID    int32
	Amount        decimal.Decimal
	Name          string
	Email         string
}

type UpdatePayment struct {
	ExternalID    string
	RawData       string
	PaidAt        string
	Status        int32
	PaymentMethod string
	Vendor        string
	Amount        decimal.Decimal
}

type GetPaginatedDonaturParams struct {
	Limit  int32
	Offset int32
	Slug   string
}

type DonaturList struct {
	ID           int32           `json:"id"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	TotalDonated decimal.Decimal `json:"total_donated"`
}
