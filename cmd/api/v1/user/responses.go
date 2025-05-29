package user

import "go-campaign.com/internal/campaign"

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

func listCampaignCollection(cs []campaign.Campaign) []ListCampaign {
	var campaigns []ListCampaign
	for _, c := range cs {
		campaigns = append(campaigns, ListCampaign{
			ID:            c.ID,
			Title:         c.Title,
			Description:   c.Description,
			TargetAmount:  c.TargetAmount,
			CurrentAmount: c.CurrentAmount,
			Slug:          c.Slug,
			StartDate:     c.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:       c.EndDate.Format("2006-01-02 15:04:05"),
			Status:        int(c.Status),
		})
	}
	return campaigns
}
