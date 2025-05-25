package validation

import (
	"database/sql"
	"strings"

	"github.com/go-playground/validator/v10"
)

func uniqueRule(db *sql.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		params := strings.Split(fl.Param(), ",") // Assuming the table name is passed as a parameter
		if len(params) < 1 {
			return false // Invalid rule, no table name provided
		}

		table := params[0]

		column := "id" // Default column to check, can be modified as needed
		if len(params) > 1 {
			column = params[1] // Use the second parameter as the column name
		}

		query := "SELECT EXISTS (SELECT 1 FROM " + table + " WHERE " + column + " = $1)"
		var exists bool
		err := db.QueryRow(query, field).Scan(&exists)
		if err != nil {
			return false // Error occurred while querying the database
		}
		return !exists // Return true if the value does not exist in the table
	}
}
