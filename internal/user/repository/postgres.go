package repository

import (
	"database/sql"
	"time"

	"go-campaign.com/internal/user"
)

type PostgresRepository struct {
	connection *sql.DB
}

func newPostgresRepository(connection *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		connection: connection,
	}
}

func (r *PostgresRepository) Create(user user.User) (user.User, error) {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	user.CreatedAt = time.Now().Format(time.RFC3339)
	user.UpdatedAt = user.CreatedAt

	err := r.connection.QueryRow(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return user, err
	}

	return user, nil
}
