package v1

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
