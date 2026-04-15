package service

import (
	"context"

	"github.com/google/uuid"
	"go-campaign.com/internal/payment/service/repository"
)

type PaymentService struct {
	repository repository.PaymentRepository
}

func NewPaymentService(repository repository.PaymentRepository) *PaymentService {
	return &PaymentService{
		repository: repository,
	}
}

func (s *PaymentService) GetDetailPayment(ctx context.Context, transactionID uuid.UUID) (*repository.DetailPayment, error) {
	return s.repository.GetDetailPayment(ctx, transactionID)
}
