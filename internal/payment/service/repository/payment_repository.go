package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PaymentRepository interface {
	GetDetailPayment(context.Context, uuid.UUID) (*DetailPayment, error)
}

type Donatur struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Donation struct {
	ID   int64  `json:"id"`
	Note string `json:"note"`
}

type Campaign struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

type DetailPayment struct {
	ID            int64           `json:"id"`
	TransactionID uuid.UUID       `json:"transaction_id"`
	Donatur       Donatur         `json:"donatur"`
	Campaign      Campaign        `json:"campaign"`
	Vendor        string          `json:"vendor"`
	Method        string          `json:"method"`
	Link          string          `json:"link"`
	Status        string          `json:"status"`
	Amount        decimal.Decimal `json:"amount"`
	PaymentDate   time.Time       `json:"payment_date"`
}
