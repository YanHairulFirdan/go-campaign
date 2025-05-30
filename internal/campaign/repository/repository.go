package repository

import (
	"database/sql"
	"fmt"

	"go-campaign.com/internal/campaign/entities"
)

type Filter struct {
	Column   string
	Value    any
	Operator string
}

type Filters []Filter

func (f Filters) ToSQL() (string, []any) {
	var conditions []string
	var args []any

	for idx, filter := range f {
		if filter.Operator == "" {
			filter.Operator = "=" // Default operator
		}
		conditions = append(conditions,
			fmt.Sprintf(
				"%s %s $%d",
				filter.Column, filter.Operator, idx+1,
			),
		)
		args = append(args, filter.Value)
	}

	if len(conditions) == 0 {
		return "", nil // No filters provided
	}

	query := "WHERE " + conditions[0]
	for _, condition := range conditions[1:] {
		query += " AND " + condition
	}

	return query, args
}

type Repository interface {
	Paginate(filters Filters, page, perPage int) ([]entities.Campaign, error)
	Create(entities.Campaign) (entities.Campaign, error)
	FindBy(column string, value any) (entities.Campaign, error)
	GetCampaignsFromUser(userID int) ([]entities.Campaign, error)
	Update(entities.Campaign) (entities.Campaign, error)
}

func NewRepository(connection *sql.DB) Repository {
	return newPostgresRepository(connection)
}
