package v1

import "mime/multipart"

type createCampaignRequest struct {
	Title        string                  `form:"title" validate:"required,min=3,max=100"`
	Description  string                  `form:"description" validate:"required,min=10,max=500"`
	Slug         string                  `form:"slug" validate:"required,min=3,max=50,unique=campaigns:slug"` // Unique slug for the campaign
	TargetAmount float32                 `form:"target_amount" validate:"required,min=0"`
	StartDate    string                  `form:"start_date" validate:"required,datetime=2006-01-02 15:04:00"`
	EndDate      string                  `form:"end_date" validate:"required,datetime=2006-01-02 15:04:00"`
	Status       int                     `form:"status" validate:"oneof=1 2,required"`                                    // 1: Draft, 2: Active
	Images       []*multipart.FileHeader `validate:"required,dive,file_max_size=1MB,file_mime_type=image/jpeg-image/png"` // Image file for the campaign
}

type updateCampaignRequest struct {
	Title        string  `form:"title" validate:"required,min=3,max=100"`
	Description  string  `form:"description" validate:"required,min=10,max=500"`
	Slug         string  `form:"slug" validate:"required,min=3,max=50,unique=campaigns:slug"` // Unique slug for the campaign
	TargetAmount float32 `form:"target_amount" validate:"required,min=0"`
	StartDate    string  `form:"start_date" validate:"required,datetime=2006-01-02 15:04:00"`
	EndDate      string  `form:"end_date" validate:"required,datetime=2006-01-02 15:04:00"`
	Status       int     `form:"status" validate:"required,oneof=1 2"` // 0: Draft, 1: Active
}
