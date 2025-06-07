package payment

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestCreateXenditTransaction(t *testing.T) {
	err := godotenv.Load("../../../../.env")

	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	x := NewXendit()
	request := InvoiceRequest{
		ExternalID: uuid.New(),
		Amount:     30000,
		Currency:   "IDR",
		UserDetail: UserDetail{
			Email:    "mamank@mail.com",
			FullName: "Mamank",
		},
		ProductDetails: []ProductDetail{
			{
				Name:     "Test Product",
				Price:    10000,
				Quantity: 1,
			},
			{
				Name:     "Test Product 2",
				Price:    20000,
				Quantity: 2,
			},
		},
	}

	res, err := x.CreateInvoice(request)

	if err != nil {
		t.Errorf("Failed to create Xendit transaction: %v", err)
		return
	}

	t.Logf("Xendit transaction created successfully: %s", res)
	fmt.Printf("Xendit transaction created successfully: %s\n", res)
	if res == "" {
		t.Error("Expected a non-empty response from Xendit, got empty string")
		return
	}
}
