package repository

import (
	"database/sql"

	"go-campaign.com/internal/user"
)

type Filter struct {
	Column   string
	Value    any
	Operator string
}

type Filters []Filter

type Repository interface {
	Create(user.User) (user.User, error)
	FindBy(column string, value any) (user.User, error)
}

func NewRepository(connection *sql.DB) Repository {
	return newPostgresRepository(connection)
}
