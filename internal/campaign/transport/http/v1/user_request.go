package v1

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validationPkg "go-campaign.com/pkg/validation"
)

type createCampaignRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Slug         string   `json:"slug"`
	TargetAmount float32  `json:"target_amount"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Status       int      `json:"status"`
	Images       []string `json:"images"`
}

func (r *createCampaignRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Title, validation.Required, validation.Length(3, 100)),
		validation.Field(&r.Description, validation.Required, validation.Length(10, 500)),
		validation.Field(&r.Slug, validation.Required, validation.Length(3, 50), validationPkg.Unique("campaigns", "slug", "", nil, "Slug already taken")),
		validation.Field(&r.TargetAmount, validation.Required, validation.Min(0.0)),
		validation.Field(&r.StartDate, validation.Required, validation.Date("2006-01-02 15:04:00")),
		validation.Field(&r.EndDate, validation.Required, validation.Date("2006-01-02 15:04:00")),
		validation.Field(&r.Status, validation.Required, validation.In(1, 2)),
		validation.Field(&r.Images, validation.Each(is.URL)),
	)
}

type updateCampaignRequest struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Slug         string   `json:"slug"`
	TargetAmount float32  `json:"target_amount"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Images       []string `json:"images"`
	Status       int      `json:"status"`
}

func (r *updateCampaignRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Title, validation.Required, validation.Length(3, 100)),
		validation.Field(&r.Description, validation.Required, validation.Length(10, 500)),
		validation.Field(&r.Slug, validation.Required, validation.Length(3, 50), validationPkg.Unique("campaigns", "slug", "id", r.ID, "Slug already taken")),
		validation.Field(&r.TargetAmount, validation.Required, validation.Min(0.0)),
		validation.Field(&r.StartDate, validation.Required, validation.Date("2006-01-02 15:04:00")),
		validation.Field(&r.EndDate, validation.Required, validation.Date("2006-01-02 15:04:00")),
		validation.Field(&r.Status, validation.Required, validation.In(1, 2)),
		validation.Field(&r.Images, validation.Each(is.URL)),
	)
}
