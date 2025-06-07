package payment

import (
	"context"
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
