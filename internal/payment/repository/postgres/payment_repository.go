package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go-campaign.com/internal/payment/repository/sqlc"
	"go-campaign.com/internal/payment/service/repository"
)

type paymentRepository struct {
	sqlc *sqlc.Queries
}

var _ repository.PaymentRepository = (*paymentRepository)(nil)

func NewPaymentRepository(q *sqlc.Queries) *paymentRepository {
	return &paymentRepository{sqlc: q}
}

func (r *paymentRepository) GetDetailPayment(ctx context.Context, transactionID uuid.UUID) (*repository.DetailPayment, error) {
	p, err := r.sqlc.GetDetailPayment(ctx, transactionID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve the payment: %w", err)
	}

	return &repository.DetailPayment{
		ID:            int64(p.ID),
		TransactionID: transactionID,
		Donatur: repository.Donatur{
			ID:   int64(p.DonaturID),
			Name: p.DonaturName,
		},
		Campaign: repository.Campaign{
			ID:          int64(p.CampaignID),
			Title:       p.CampaignTitle,
			Description: p.CampaignDescription.String,
			Creator:     p.Creator,
		},
		Vendor:      p.Vendor.String,
		Method:      p.Method.String,
		Link:        p.Link.String,
		Status:      "p.Status",
		Amount:      p.Amount,
		PaymentDate: p.PaymentDate.Time,
	}, nil
}
