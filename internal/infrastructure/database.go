package infrastructure

import (
	"database/sql"
	"fmt"

	"go-campaign.com/internal/config"
)

func InitDatabaseConnection(c *config.Config) (*sql.DB, error) {
	dbConnectionString := c.Database.URL

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))
	}

	return db, nil
}
