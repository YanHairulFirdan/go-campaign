package repository

import (
	"database/sql"
	"errors"
	"time"

	"go-campaign.com/internal/campaign"
)

type PostgresRepository struct {
	connection *sql.DB
}

func newPostgresRepository(connection *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		connection: connection,
	}
}

func (r *PostgresRepository) Create(c campaign.Campaign) (campaign.Campaign, error) {
	query := `INSERT INTO 
		campaigns (user_id, title, description, slug, target_amount, current_amount, start_date, end_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`
	createdAt := time.Now()
	updatedAt := time.Now()
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt
	err := r.connection.QueryRow(query,
		c.UserID,
		c.Title,
		c.Description,
		c.Slug,
		c.TargetAmount,
		c.CurrentAmount,
		c.StartDate,
		c.EndDate,
		c.Status,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return campaign.Campaign{}, err // Return error if insertion fails
	}

	return c, nil // Return the created campaign
}

func (r *PostgresRepository) FindBy(column string, value any) (campaign.Campaign, error) {
	query := `SELECT id, user_id, title, description, slug, target_amount, current_amount, start_date, end_date, status, created_at, updated_at, deleted_at
		FROM campaigns WHERE ` + column + ` = $1`

	var c campaign.Campaign
	err := r.connection.QueryRow(query, value).Scan(
		&c.ID,
		&c.UserID,
		&c.Title,
		&c.Description,
		&c.Slug,
		&c.TargetAmount,
		&c.CurrentAmount,
		&c.StartDate,
		&c.EndDate,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return campaign.Campaign{}, errors.New("campaign not found")
		}
		return campaign.Campaign{}, err // Return error if query fails
	}

	return c, nil
}

func (r *PostgresRepository) GetCampaignsFromUser(userID int) ([]campaign.Campaign, error) {
	query := `SELECT id, user_id, title, description, slug, target_amount, current_amount, start_date, end_date, status, created_at, updated_at, deleted_at
		FROM campaigns WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.connection.Query(query, userID)
	if err != nil {
		return nil, err // Return error if query fails
	}
	defer rows.Close()

	var campaigns []campaign.Campaign
	for rows.Next() {
		var c campaign.Campaign
		err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Title,
			&c.Description,
			&c.Slug,
			&c.TargetAmount,
			&c.CurrentAmount,
			&c.StartDate,
			&c.EndDate,
			&c.Status,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)
		if err != nil {
			return nil, err // Return error if scanning fails
		}
		campaigns = append(campaigns, c)
	}

	return campaigns, nil // Return the list of campaigns
}

func (r *PostgresRepository) Update(c campaign.Campaign) (campaign.Campaign, error) {
	query := `
	UPDATE campaigns SET
	 title = $1,
	 slug = $2,
	 description = $3,
	 target_amount = $4,
	 current_amount = $5,
	 start_date = $6,
	 end_date = $7,
	 status = $8,
	 updated_at = $9,
	 deleted_at = $10
	WHERE id = $11
	RETURNING id, user_id, title, description, slug, target_amount, current_amount, start_date, end_date, status, created_at, updated_at, deleted_at`
	updatedAt := time.Now()
	c.UpdatedAt = updatedAt
	err := r.connection.QueryRow(query,
		c.Title,
		c.Slug,
		c.Description,
		c.TargetAmount,
		c.CurrentAmount,
		c.StartDate,
		c.EndDate,
		c.Status,
		c.UpdatedAt,
		c.DeletedAt,
		c.ID,
	).Scan(
		&c.ID,
		&c.UserID,
		&c.Title,
		&c.Description,
		&c.Slug,
		&c.TargetAmount,
		&c.CurrentAmount,
		&c.StartDate,
		&c.EndDate,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	if err != nil {
		return campaign.Campaign{}, err // Return error if update fails
	}

	return c, nil // Return the updated campaign
}
