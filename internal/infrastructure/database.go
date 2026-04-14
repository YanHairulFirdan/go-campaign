package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-campaign.com/internal/config"
)

func InitDatabaseConnection(c *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", c.Database.URL)

	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()

		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
