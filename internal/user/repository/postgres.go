package repository

import (
	"database/sql"
	"errors"
	"time"

	"go-campaign.com/internal/user/entities"
)

type PostgresRepository struct {
	connection *sql.DB
}

func newPostgresRepository(connection *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		connection: connection,
	}
}

func (r *PostgresRepository) Create(user entities.User) (entities.User, error) {
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

func (r *PostgresRepository) FindBy(column string, value any) (entities.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at 
	          FROM users WHERE ` + column + ` = $1`

	var u entities.User

	err := r.connection.QueryRow(query, value).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errors.New("user not found") // No user found
		}
		return entities.User{}, err // Other error
	}

	return u, nil
}
