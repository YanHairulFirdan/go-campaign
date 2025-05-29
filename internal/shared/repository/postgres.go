package repository

import (
	"database/sql"
	"fmt"
)

type postgresStore struct {
	db *sql.DB
}

func newPostgresStore(db *sql.DB) *postgresStore {
	return &postgresStore{
		db: db,
	}
}

func (s *postgresStore) IsUnique(table, column string, value any) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM " + table + " WHERE " + column + " = $1)"
	var exists bool
	err := s.db.QueryRow(query, value).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *postgresStore) IsUniqueWithCondition(
	table,
	column,
	conditionColumn string,
	value,
	conditionValue any,
) (bool, error) {
	query := fmt.Sprintf(
		"SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $1 AND %s != $2)",
		table, column, conditionColumn,
	)
	var exists bool
	err := s.db.QueryRow(query, value, conditionValue).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
