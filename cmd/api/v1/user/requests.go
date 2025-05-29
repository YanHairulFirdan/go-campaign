package user

type createCampaignRequest struct {
	Title        string  `json:"title" validate:"required,min=3,max=100"`
	Description  string  `json:"description" validate:"required,min=10,max=500"`
	Slug         string  `json:"slug" validate:"required,min=3,max=50,unique=campaigns:slug"` // Unique slug for the campaign
	TargetAmount float32 `json:"target_amount" validate:"required,min=0"`
	StartDate    string  `json:"start_date" validate:"required,datetime=2006-01-02 15:04:00"`
	EndDate      string  `json:"end_date" validate:"required,datetime=2006-01-02 15:04:00"`
	Status       int     `json:"status" validate:"required,oneof=0 1"` // 0: Draft, 1: Active
}

type updateCampaignRequest struct {
	Title        string  `json:"title" validate:"required,min=3,max=100"`
	Description  string  `json:"description" validate:"required,min=10,max=500"`
	Slug         string  `json:"slug" validate:"required,min=3,max=50,unique=campaigns:slug"` // Unique slug for the campaign
	TargetAmount float32 `json:"target_amount" validate:"required,min=0"`
	StartDate    string  `json:"start_date" validate:"required,datetime=2006-01-02 15:04:00"`
	EndDate      string  `json:"end_date" validate:"required,datetime=2006-01-02 15:04:00"`
	Status       int     `json:"status" validate:"required,oneof=0 1"` // 0: Draft, 1: Active
}
