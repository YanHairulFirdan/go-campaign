package v1

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ListCampaign struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	TargetAmount  float32 `json:"target_amount"`
	CurrentAmount float32 `json:"current_amount"`
	Slug          string  `json:"slug"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Status        int     `json:"status"` // 0: Draft, 1: Active, 2: Completed, 3: Cancelled
}

type DonationRequest struct {
	Name   string  `json:"name" validate:"required,min=3,max=100"`
	Email  string  `json:"email" validate:"required,email"`
	Amount float32 `json:"amount" validate:"required,min=1"`
	Note   string  `json:"note" validate:"omitempty,max=500"`
}

func (r *DonationRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&r.Email, validation.Required, validation.Length(5, 100), is.Email),
		validation.Field(&r.Amount, validation.Required, validation.Min(1.0)),
		validation.Field(&r.Note, validation.Length(0, 500)),
	)
}
