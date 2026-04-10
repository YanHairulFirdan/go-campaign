package payment

import "github.com/shopspring/decimal"

type InvoiceStatus string

const (
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusExpired InvoiceStatus = "EXPIRED"
)

type InvoiceWebhookResponseDto struct {
	ID                 string
	ExternalID         string
	Status             InvoiceStatus
	Amount             decimal.Decimal
	PayerEmail         string
	Description        string
	PaidAmount         decimal.Decimal
	Updated            string
	Created            string
	Currency           string
	PaidAt             string
	PaymentMethod      string
	PaymentChannel     string
	PaymentDetails     map[string]interface{}
	PaymentID          string
	SuccessRedirectURL string
	FailureRedirectURL string
	Items              []InvoiceItem
}

type InvoiceItem struct {
	Name     string
	Quantity float32
	Price    float32
	Category string
}
