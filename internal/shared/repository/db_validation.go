package repository

import "database/sql"

type DatabaseValidationRepository interface {
	IsUnique(table, column string, value any) (bool, error)
	IsUniqueWithCondition(table, column, conditionColumn string, value, conditionValue any) (bool, error)
}

func NewDatabaseValidationRepository(db *sql.DB) DatabaseValidationRepository {
	return newPostgresStore(db)
}
