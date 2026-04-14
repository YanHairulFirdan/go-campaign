package services

import (
	"context"
	"fmt"

	"go-campaign.com/internal/campaign/services/repository"
	"go-campaign.com/internal/shared/http/request"
	"go-campaign.com/internal/shared/services/payment"
)

type CampaignService struct {
	p                  payment.PaymentGateway
	donationRepository repository.DonationRepository
	campaignRepository repository.CampaignRepository
}

func NewCampaignService(
	p payment.PaymentGateway,
	donationRepository repository.DonationRepository,
	campaignRepository repository.CampaignRepository,
) *CampaignService {
	return &CampaignService{
		p:                  p,
		donationRepository: donationRepository,
		campaignRepository: campaignRepository,
	}
}

func (s *CampaignService) GetCampaigns(ctx context.Context, offset, limit int32) ([]repository.CampaignList, int, error) {
	campaigns, err := s.campaignRepository.GetPaginatedCampaigns(ctx, request.PaginationRequest{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get campaigns: %w", err)
	}

	totalCount, err := s.campaignRepository.GetTotalCampaign(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total campaigns count: %w", err)
	}

	return campaigns, int(totalCount), nil
}

func (s *CampaignService) GetCampaignBySlug(ctx context.Context, slug string) (*repository.DetailCampaign, error) {
	campaign, err := s.campaignRepository.GetCampaignBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return campaign, nil
}

func (s *CampaignService) Donate(ctx context.Context, request DonationRequest) (string, error) {
	var url string

	invoiceIntent, err := s.donationRepository.CreateDonationIntent(ctx, repository.CreateDonationIntentParams(request))

	if err != nil {
		return "", fmt.Errorf("failed to create donation intent: %w", err)
	}

	url, err = s.p.CreateInvoice(payment.InvoiceRequest{
		ExternalID: invoiceIntent.TransactionID,
		Amount:     invoiceIntent.Amount.InexactFloat64(),
		Currency:   "IDR",
		UserDetail: payment.UserDetail{
			Email:    request.Email,
			FullName: request.Name,
		},
		ProductDetails: []payment.ProductDetail{
			{
				Name:     "Donation to campaign",
				Price:    request.Amount.Abs().InexactFloat64(),
				Quantity: 1,
			},
		},
	})

	if err != nil {
		invoiceErr := s.donationRepository.MarkInvoiceFailed(ctx, invoiceIntent.PaymentID)

		if invoiceErr != nil {
			return "", fmt.Errorf("failed to mark invoice as failed: %w", invoiceErr)
		}

		return "", fmt.Errorf("failed to create payment invoice: %w", err)
	}

	invoiceErr := s.donationRepository.MarkInvoiceCreated(ctx, invoiceIntent.PaymentID, url, "xendit")

	if invoiceErr != nil {
		return "", fmt.Errorf("failed to mark invoice as created: %w", invoiceErr)
	}

	return url, nil
}

func (s *CampaignService) UpdatePaymentFromCallback(ctx context.Context, webhookEvent *payment.PaymentCallback) error {
	status, exists := payment.MapPaymentStatus[webhookEvent.Status]

	if !exists {
		return fmt.Errorf("payment status is invalid: %w", webhookEvent.Status)
	}

	err := s.donationRepository.UpdateDonationPaymentFromWebhook(ctx, repository.UpdatePayment{
		ExternalID:    webhookEvent.ExternalID,
		RawData:       webhookEvent.RawData,
		PaidAt:        webhookEvent.PaidAt,
		PaymentMethod: webhookEvent.PaymentMethod,
		Status:        status,
	})

	return err
}

func (s *CampaignService) GetDonatur(ctx context.Context, req GetDonaturListRequest) ([]repository.DonaturList, int, error) {
	donaturs, err := s.donationRepository.GetPaginatedDonatur(ctx, repository.GetPaginatedDonaturParams{
		Offset: req.Offset,
		Limit:  req.Limit,
		Slug:   req.Slug,
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated donaturs: %w", err)
	}

	totalCount, err := s.donationRepository.GetTotalPaidDonatur(ctx, req.Slug)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total paid donaturs: %w", err)
	}

	return donaturs, int(totalCount), nil
}
