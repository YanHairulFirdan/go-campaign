package validation

import (
	"errors"
)

func Unique(table string, field string, excludeColumn string, excludeValue any, message string) dbUniqueRule {
	if message == "" {
		message = field + " already taken"
	}

	return dbUniqueRule{
		table:         table,
		field:         field,
		excludeColumn: excludeColumn,
		excludeValue:  excludeValue,
		message:       message,
	}
}

type dbUniqueRule struct {
	table         string
	field         string
	excludeColumn string
	excludeValue  any
	message       string
}

func (r dbUniqueRule) Validate(value interface{}) error {
	var exists bool
	var err error

	if r.excludeColumn != "" {
		exists, err = databaseRepository.IsUniqueWithCondition(r.table, r.field, r.excludeColumn, value, r.excludeValue)
	} else {
		exists, err = databaseRepository.IsUnique(r.table, r.field, value)
	}

	if err != nil {
		return errors.New("error occurred while checking uniqueness")
	}

	if exists {
		return errors.New(r.message)
	}
	return nil
}
