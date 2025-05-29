package validation

import (
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

func uniqueRule(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	fieldName := fl.FieldName()

	params := strings.Split(fl.Param(), ":") // Assuming the table name is passed as a parameter
	if len(params) < 1 {
		return false // Invalid rule, no table name provided
	}

	table := params[0]

	column := "id" // Default column to check, can be modified as needed

	if len(params) > 1 {
		column = params[1] // Use the second parameter as the column name
	}

	var exists bool
	var err error

	if exception, exceptionExists := validationFieldExceptions[fieldName]; exceptionExists {
		exists, err = databaseRepository.IsUniqueWithCondition(table, column, exception.Column, value, exception.Value)

		if err != nil {
			log.Printf("Error checking uniqueness with exception: %v\n", err)
			return false // Error occurred while checking uniqueness
		}
	} else {
		exists, err = databaseRepository.IsUnique(table, column, value)
		if err != nil {
			log.Printf("Error checking uniqueness: %v\n", err)
			return false // Error occurred while checking uniqueness
		}
	}

	return !exists // Return true if the value does not exist in the table
}
