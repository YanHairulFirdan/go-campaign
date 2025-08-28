package payment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserDetail struct {
	Email    string
	FullName string
}

type ProductDetail struct {
	Name        string
	Description string
	Price       float64 // Price in the smallest currency unit (e.g., cents for USD)
	Quantity    int
}

type InvoiceRequest struct {
	ExternalID     uuid.UUID // Unique identifier for the invoice, can be used to track the transaction
	Amount         float64
	Currency       string // e.g., "USD", "IDR"
	UserDetail     UserDetail
	ProductDetails []ProductDetail
}

type PaymentStatus string

const (
	PaymentStatusPending      PaymentStatus = "PENDING"
	PaymentStatusProcessing   PaymentStatus = "PROCESSING"
	PaymentStatusSuccess      PaymentStatus = "SUCCESS"
	PaymentStatusFailed       PaymentStatus = "FAILED"
	PaymentStatusExpired      PaymentStatus = "EXPIRED"
	DonationPaymentStatusPaid PaymentStatus = "PAID"
)

type PaymentCallback struct {
	ExternalID    string
	RawData       string
	PaidAt        string
	Status        PaymentStatus
	PaymentMethod string
}

type PaymentStatusMap map[PaymentStatus]int32

var MapPaymentStatus = map[PaymentStatus]int32{
	PaymentStatusPending:      0,
	PaymentStatusProcessing:   1,
	PaymentStatusSuccess:      2,
	PaymentStatusFailed:       3,
	PaymentStatusExpired:      4,
	DonationPaymentStatusPaid: 5,
}

type PaymentGateway interface {
	CreateInvoice(request InvoiceRequest) (string, error) // Returns a payment URL or ID
	ParseCallbackResponse(c *fiber.Ctx) (*PaymentCallback, error)
}

func New() PaymentGateway {
	return NewXendit()
}
