package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/common"
	"github.com/xendit/xendit-go/v7/invoice"
)

type xenditPaymentGateway struct {
	client *xendit.APIClient // Xendit client for API interactions
}

type InvoiceStatus string

const (
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusExpired InvoiceStatus = "EXPIRED"
)

type XenditInvoiceItem struct {
	Name     string  `json:"name"`
	Quantity float32 `json:"quantity"`
	Price    float32 `json:"price"`
	Category string  `json:"category,omitempty"` // Optional category field
}

type XenditInvoiceWebhookResponse struct {
	ID                 string                 `json:"id"`
	ExternalID         string                 `json:"external_id"`
	UserID             string                 `json:"user_id"`
	Status             InvoiceStatus          `json:"status"`
	MerchantName       string                 `json:"merchant_name"`
	Amount             float64                `json:"amount"`
	PayerEmail         string                 `json:"payer_email"`
	Description        string                 `json:"description"`
	PaidAmount         float64                `json:"paid_amount"`
	Updated            string                 `json:"updated"`
	Created            string                 `json:"created"`
	Currency           string                 `json:"currency"`
	PaidAt             string                 `json:"paid_at"`
	PaymentMethod      string                 `json:"payment_method"`
	PaymentChannel     string                 `json:"payment_channel"`
	PaymentDetails     map[string]interface{} `json:"payment_details"`
	PaymentID          string                 `json:"payment_id"`
	SuccessRedirectURL string                 `json:"success_redirect_url"`
	FailureRedirectURL string                 `json:"failure_redirect_url"`
	Items              []XenditInvoiceItem    `json:"items"`
}

func (r *XenditInvoiceWebhookResponse) ToJson() (string, error) {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Xendit invoice webhook response: %w", err)
	}
	return string(jsonData), nil
}

func NewXendit() *xenditPaymentGateway {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	if secretKey == "" {
		panic("XENDIT_SECRET_KEY environment variable is not set")
	}

	return &xenditPaymentGateway{
		client: xendit.NewClient(secretKey),
	}
}

// CreateRequest creates a payment request using Xendit.
func (x *xenditPaymentGateway) CreateInvoice(request InvoiceRequest) (string, error) {
	var totalAmount float64

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest("some-invoice-id", request.Amount)
	createInvoiceRequest.Currency = &request.Currency
	createInvoiceRequest.ExternalId = request.ExternalID.String()
	createInvoiceRequest.Customer = &invoice.CustomerObject{
		Email: *invoice.NewNullableString(&request.UserDetail.Email),
	}

	for _, product := range request.ProductDetails {
		totalAmount += product.Price * float64(product.Quantity)
		createInvoiceRequest.Items = append(createInvoiceRequest.Items, invoice.InvoiceItem{
			Name:     product.Name,
			Price:    float32(product.Price),
			Quantity: float32(product.Quantity),
		})
	}
	createInvoiceRequest.Amount = totalAmount

	resp, r, err := x.client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if err != nil {
		var sdkErr *common.XenditSdkError

		if errors.As(err, &sdkErr) {
			return "", fmt.Errorf("xendit SDK error: %s, HTTP status: %d, Response: %v", sdkErr.Error(), r.StatusCode, sdkErr)
		}

		return "", fmt.Errorf("failed to create Xendit invoice: %w, HTTP status: %d, Response: %v", err, r.StatusCode, r.Body)
	}

	return resp.InvoiceUrl, nil
}
