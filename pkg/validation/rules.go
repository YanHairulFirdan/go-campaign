package validation

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

func uniqueRule(db *sql.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
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

		query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $1)", table, column)

		if exception, exceptionExists := validationFieldExceptions[fieldName]; exceptionExists {
			// If there is an exception for this field, modify the query accordingly
			query = fmt.Sprintf(
				"SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $1 AND %s != $2)",
				table,
				column,
				exception.Column,
			)

			// Assuming exception.Value contains the ID to exclude from the uniqueness check
			err := db.QueryRow(query, value, exception.Value).Scan(&exists)
			if err != nil {
				return false // Error occurred while querying the database
			}
			return !exists // Return true if the value does not exist in the table
		}

		err := db.QueryRow(query, value).Scan(&exists)

		log.Println("Executing query:", query, "with value:", value)
		log.Println("Query result:", exists)
		if err != nil {
			log.Printf("Error executing query: %v\n", err)
			return false // Error occurred while querying the database
		}
		return !exists // Return true if the value does not exist in the table
	}
}
