package entities

import "time"

type Status int

const (
	StatusDraft     Status = 1
	StatusActive    Status = 2
	StatusCompleted Status = 3
	StatusCancelled Status = 4
)

type Campaign struct {
	ID            int
	UserID        int
	Title         string
	Description   string
	Slug          string
	TargetAmount  float32
	CurrentAmount float32
	StartDate     time.Time
	EndDate       time.Time
	Status        Status
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
