package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sqlc-dev/pqtype"
	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/campaign/services/repository"
)

type DonationRepository struct {
	db   *sql.DB
	sqlc *sqlc.Queries
}

func NewDonationRepository(db *sql.DB, sqlc *sqlc.Queries) *DonationRepository {
	return &DonationRepository{
		db:   db,
		sqlc: sqlc,
	}
}

var _ repository.DonationRepository = (*DonationRepository)(nil)

func (r *DonationRepository) CreateDonationIntent(ctx context.Context, req repository.CreateDonationIntentParams) (*repository.DonationIntent, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("donation amount must not zero")
	}

	if req.Amount.Exponent() < -2 {
		return nil, fmt.Errorf("amount must have at most 2 decimal places")
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	qtx := r.sqlc.WithTx(tx)

	donatur, err := qtx.CreateDonatur(ctx, sqlc.CreateDonaturParams{
		Name: req.Name,
		Email: sql.NullString{
			String: req.Email,
			Valid:  req.Email != "",
		},
		UserID:     req.UserID,
		CampaignID: req.CampaignID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create donatur: %w", err)
	}

	var note = ""

	if req.Note != nil {
		note = *req.Note
	}

	donation, err := qtx.CreateDonation(ctx, sqlc.CreateDonationParams{
		DonaturID:  donatur.ID,
		CampaignID: req.CampaignID,
		Amount:     req.Amount,
		Note: sql.NullString{
			String: note,
			Valid:  req.Note != nil,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create donation: %w", err)
	}

	transactionID := uuid.New()

	payment, err := qtx.CreatePayment(ctx, sqlc.CreatePaymentParams{
		TransactionID: transactionID,
		DonationID:    donation.ID,
		DonaturID:     donatur.ID,
		CampaignID:    req.CampaignID,
		Amount:        req.Amount,
		Status:        int32(sqlc.DonationPaymentStatusPending),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &repository.DonationIntent{
		PaymentID:     payment.ID,
		TransactionID: payment.TransactionID,
		CampaignID:    donation.CampaignID,
		Amount:        req.Amount,
		Name:          donatur.Name,
		Email:         donatur.Email.String,
	}, nil
}

func (r *DonationRepository) MarkInvoiceCreated(ctx context.Context, paymentID int32, link, vendor string) error {
	return r.sqlc.MarkPaymentInvoiceCreated(ctx, sqlc.MarkPaymentInvoiceCreatedParams{
		Link: sql.NullString{
			String: link,
			Valid:  link != "",
		},
		Vendor: sql.NullString{
			String: vendor,
			Valid:  vendor != "",
		},
		Status: int32(sqlc.DonationPaymentStatusProcessing),
		ID:     paymentID,
	})
}

func (r *DonationRepository) MarkInvoiceFailed(ctx context.Context, paymentID int32) error {
	return r.sqlc.UpdatePaymentStatus(ctx, sqlc.UpdatePaymentStatusParams{
		Status: int32(sqlc.DonationPaymentStatusRetry),
		ID:     paymentID,
	})
}

func (r *DonationRepository) UpdateDonationPaymentFromWebhook(ctx context.Context, req repository.UpdatePayment) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to start the database transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	payment, err := r.sqlc.FindAndLockPaymentForUpdate(ctx, uuid.MustParse(req.ExternalID))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("the payment record not found")
		}

		return fmt.Errorf("failed to find the payment: %w", err)
	}

	updateParams := sqlc.UpdatePaymentFromWebhookCallbackParams{
		Status: req.Status,
		ID:     payment.ID,
		Vendor: sql.NullString{
			String: req.Vendor,
			Valid:  req.Vendor != "",
		},
		Method: sql.NullString{
			String: req.PaymentMethod,
			Valid:  req.PaymentMethod != "",
		},
		Response: pqtype.NullRawMessage{
			RawMessage: []byte(req.RawData),
			Valid:      true,
		},
		PaymentDate: sql.NullTime{
			Time: func() time.Time {
				parsedTime, err := time.Parse("2006-01-02T15:04:05Z", req.PaidAt)

				if err != nil {
					return time.Time{}
				}

				return parsedTime
			}(),
			Valid: req.PaidAt != "",
		},
	}

	_, err = r.sqlc.UpdatePaymentFromWebhookCallback(ctx, updateParams)

	if err != nil {
		return fmt.Errorf("failed to update payment status")
	}

	if sqlc.DonationPaymentStatus(req.Status) != sqlc.DonationPaymentStatusPaid {
		return tx.Commit()
	}

	campaign, err := r.sqlc.FindCampaignByIdForUpdate(ctx, payment.CampaignID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("campaign not found")
		}

		return fmt.Errorf("failed to retrieve the campaign: %w", err)
	}

	err = r.sqlc.IncreaseCampaignCurrentAmount(ctx, sqlc.IncreaseCampaignCurrentAmountParams{
		ID:     campaign.ID,
		Amount: req.Amount,
	})

	if err != nil {
		return fmt.Errorf("failed to increase the campaign's current_amount: %w", err)
	}

	return tx.Commit()
}

func (r *DonationRepository) GetPaginatedDonatur(ctx context.Context, req repository.GetPaginatedDonaturParams) ([]repository.DonaturList, error) {
	donaturs, err := r.sqlc.GetPaginatedDonaturs(ctx, sqlc.GetPaginatedDonatursParams{
		Slug:   req.Slug,
		Limit:  req.Limit,
		Offset: req.Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve donatur list: %w", err)
	}

	var donaturList []repository.DonaturList

	for _, donatur := range donaturs {
		donaturList = append(donaturList, repository.DonaturList(donatur))
	}

	return donaturList, nil
}

func (r *DonationRepository) GetTotalPaidDonatur(ctx context.Context, slug string) (int64, error) {
	total, err := r.sqlc.GetCampaignTotalPaidDonaturs(ctx, slug)

	if err != nil {
		return 0, err
	}

	return total, nil
}
