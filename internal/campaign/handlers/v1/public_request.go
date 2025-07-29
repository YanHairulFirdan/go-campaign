package v1

type DonationRequest struct {
	Name   string  `json:"name" validate:"required,min=3,max=100"`
	Email  string  `json:"email" validate:"required,email"`
	Amount float32 `json:"amount" validate:"required,min=1"`
	Note   string  `json:"note" validate:"omitempty,max=500"`
}
