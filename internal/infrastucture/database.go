package infrastucture

import (
	"database/sql"
	"fmt"
	"os"
)

func InitDatabaseConnection() (*sql.DB, error) {
	dbConnectionString := os.Getenv("DATABASE_CONNECTION")

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database: %v", err))
	}

	return db, nil
}
