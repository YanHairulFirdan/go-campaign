package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go-campaign.com/internal/campaign/repository/sqlc"
	"go-campaign.com/internal/shared/repository"
	"go-campaign.com/internal/shared/services/payment"
)

type CampaignService struct {
	q       *sqlc.Queries
	txStore *repository.TransactionStore
	p       payment.PaymentGateway
}

func NewCampaignService(q *sqlc.Queries, txStore *repository.TransactionStore, p payment.PaymentGateway) *CampaignService {
	return &CampaignService{
		q:       q,
		txStore: txStore,
		p:       p,
	}
}

func (s *CampaignService) GetCampaigns(ctx context.Context, page, perPage int32) ([]sqlc.GetCampaignsRow, int, error) {
	campaigns, err := s.q.GetCampaigns(ctx, sqlc.GetCampaignsParams{
		Limit:  perPage,
		Offset: (page - 1) * perPage,
	})

	log.Printf("GetCampaigns: page=%d, perPage=%d, total campaigns=%d", page, perPage, len(campaigns))

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get campaigns: %w", err)
	}

	totalCount, err := s.q.GetTotalCampaigns(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total campaigns count: %w", err)
	}

	return campaigns, int(totalCount), nil
}

func (s *CampaignService) GetCampaignBySlug(ctx context.Context, slug string) (*sqlc.GetCampaignBySlugRow, error) {
	campaign, err := s.q.GetCampaignBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("campaign with slug %s not found", slug)
		}

		return nil, err
	}

	return &campaign, nil
}

func (s *CampaignService) FindCampaignsBySlugForUpdate(ctx context.Context, slug string) (*sqlc.FindCampaignsBySlugForUpdateRow, error) {
	campaign, err := s.q.FindCampaignsBySlugForUpdate(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("campaign with slug %s not found", slug)
		}

		return nil, fmt.Errorf("failed to find campaign by slug for update: %w", err)
	}

	return &campaign, nil

}

func (s *CampaignService) Donate(ctx context.Context, request DonationRequest) (string, error) {
	var url string
	err := s.txStore.ExecTx(func() error {
		donatur, err := s.q.CreateDonatur(ctx, sqlc.CreateDonaturParams{
			Name: request.Name,
			Email: sql.NullString{
				String: request.Email,
				Valid:  request.Email != "",
			},
			UserID:     request.UserID,
			CampaignID: request.CampaignID,
		})

		if err != nil {
			return fmt.Errorf("failed to create donatur: %w", err)
		}

		donation, err := s.q.CreateDonation(ctx, sqlc.CreateDonationParams{
			DonaturID:  donatur.ID,
			CampaignID: request.CampaignID,
			Amount:     fmt.Sprintf("%.2f", request.Amount),
			Note: sql.NullString{
				String: *request.Note,
				Valid:  request.Note != nil,
			},
		})

		if err != nil {
			return fmt.Errorf("failed to create donation: %w", err)
		}

		transactionID := uuid.New()
		invoiceReq := payment.InvoiceRequest{
			ExternalID: transactionID,
			Amount:     float64(request.Amount),
			Currency:   "IDR",
			UserDetail: payment.UserDetail{
				Email:    request.Email,
				FullName: request.Name,
			},
			ProductDetails: []payment.ProductDetail{
				{Name: "Donation to campaign", Price: float64(request.Amount), Quantity: 1},
			},
		}

		url, err = s.p.CreateInvoice(invoiceReq)

		if err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}

		_, err = s.q.CreatePayment(ctx, sqlc.CreatePaymentParams{
			TransactionID: transactionID,
			DonationID:    donation.ID,
			DonaturID:     donatur.ID,
			CampaignID:    request.CampaignID,
			Amount:        fmt.Sprintf("%.2f", request.Amount),
			Status:        int32(sqlc.DonationPaymentStatusPending),
		})

		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("transaction failed: %w", err)
	}

	return url, nil
}

func (s *CampaignService) UpdatePaymentFromCallback(ctx context.Context, webhookEvent payment.XenditInvoiceWebhookResponse) error {
	err := s.txStore.ExecTx(func() error {
		p, err := s.q.GetPaymentByTransactionId(ctx, uuid.MustParse(webhookEvent.ExternalID))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("payment not found for transaction ID: %s", webhookEvent.ExternalID)
			}
			return fmt.Errorf("failed to get payment by transaction ID: %w", err)
		}

		rawResponse, err := webhookEvent.ToJson()

		if err != nil {
			return fmt.Errorf("failed to convert webhook event to JSON: %w", err)
		}

		updateParams := sqlc.UpdatePaymentFromCallbackParams{
			Status: int32(sqlc.DonationPaymentStatusPaid),
			ID:     p.ID,
			Vendor: sql.NullString{
				String: "xendit",
				Valid:  true,
			},
			Method: sql.NullString{
				String: webhookEvent.PaymentMethod,
				Valid:  true,
			},
			Response: pqtype.NullRawMessage{
				RawMessage: []byte(rawResponse),
				Valid:      true,
			},
			PaymentDate: sql.NullTime{
				Time: func() time.Time {
					parsedTime, err := time.Parse("2006-01-02T15:04:05Z", webhookEvent.PaidAt)
					if err != nil {
						return time.Time{}
					}
					return parsedTime
				}(),
				Valid: webhookEvent.PaidAt != "",
			}}

		if webhookEvent.Status != payment.InvoiceStatusPaid {
			updateParams.Status = int32(sqlc.DonationPaymentStatusExpired)

			_, err = s.q.UpdatePaymentFromCallback(ctx, updateParams)

			if err != nil {
				return fmt.Errorf("failed to update payment status to expired: %w", err)
			}
		}

		_, err = s.q.UpdatePaymentFromCallback(ctx, updateParams)

		if err != nil {
			return fmt.Errorf("failed to update payment from callback: %w", err)
		}

		cp, err := s.q.FindCampaignByIdForUpdate(ctx, p.CampaignID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("campaign with ID %d not found", p.CampaignID)
			}
			return fmt.Errorf("failed to find campaign by ID for update: %w", err)
		}

		err = s.q.Donate(ctx, sqlc.DonateParams{
			ID:     cp.ID,
			Amount: p.Amount,
		})

		if err != nil {
			return fmt.Errorf("failed to record donation: %w", err)
		}

		return nil
	})

	return err
}

func (s *CampaignService) GetDonatur(ctx context.Context, request GetDonaturListRequest) ([]sqlc.GetPaginatedDonatursRow, int, error) {
	donaturs, err := s.q.GetPaginatedDonaturs(ctx, sqlc.GetPaginatedDonatursParams{
		Slug:   request.Slug,
		Limit:  request.Limit,
		Offset: request.Offset,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated donaturs: %w", err)
	}

	totalCount, err := s.q.GetCampaignTotalPaidDonaturs(ctx, request.Slug)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total paid donaturs: %w", err)
	}

	return donaturs, int(totalCount), nil
}
