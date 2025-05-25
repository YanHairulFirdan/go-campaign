package repository

import (
	"database/sql"

	"go-campaign.com/internal/user"
)

type Repository interface {
	Create(user.User) (user.User, error)
}

func NewRepository(connection *sql.DB) Repository {
	return newPostgresRepository(connection)
}
