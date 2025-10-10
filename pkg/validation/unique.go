package validation

import (
	"errors"
	"log"
)

func Unique(table string, field string, excludeColumn string, excludeValue any, message string) dbUniqueRule {
	if message == "" {
		message = field + " already taken"
	}

	log.Println("Creating Unique rule for table:", table, "field:", field, "excludeColumn:", excludeColumn, "excludeValue:", excludeValue)
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
	log.Println("Validating uniqueness for ", r.table, ".", r.field, " with value: ", value)
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
