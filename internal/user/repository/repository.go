package repository

import (
	"database/sql"

	"go-campaign.com/internal/user/entities"
)

type Filter struct {
	Column   string
	Value    any
	Operator string
}

type Filters []Filter

type Repository interface {
	Create(entities.User) (entities.User, error)
	FindBy(column string, value any) (entities.User, error)
}

func NewRepository(connection *sql.DB) Repository {
	return newPostgresRepository(connection)
}
