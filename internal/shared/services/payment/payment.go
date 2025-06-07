package payment

import "github.com/google/uuid"

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

type PaymentGateway interface {
	CreateInvoice(request InvoiceRequest) (string, error) // Returns a payment URL or ID
}
