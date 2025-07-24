package validation

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
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

func fileRequiredRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)

	return ok && file != nil && file.Size > 0 // Check if the file is present and has a size greater than 0
}

func fileMaxSizeRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return false // Not a valid file header
	}

	maxSizeParam := fl.Param() // Get the max size from the validation tag parameter

	size, err := extractFileSize(maxSizeParam) // Extract size and unit type

	if err != nil {
		log.Printf("Invalid max size parameter: %v\n", err)
		return false // Invalid max size parameter
	}

	return file.Size <= size
}

func extractFileSize(fileSize string) (int64, error) {
	unitType := strings.ToUpper(fileSize[len(fileSize)-2:]) // Get the unit type (MB, KB, B)
	multipliers := map[string]int64{
		"MB": 1024 * 1024,
		"KB": 1024,
		"B":  1,
	}
	size := fileSize[:len(fileSize)-2] // Remove the unit type
	sizeParsed, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		log.Printf("Invalid file size parameter: %v\n", err)
		return 0, err // Invalid file size parameter
	}
	if multiplier, exists := multipliers[unitType]; exists {
		return sizeParsed * multiplier, nil // Convert to bytes
	} else {
		log.Printf("Unsupported unit type: %s\n", unitType)
		return 0, errors.New("unsupported unit type") // Unsupported unit type
	}
}

// min size rule
func fileMinSizeRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return false // Not a valid file header
	}

	minSizeParam := fl.Param() // Get the min size from the validation tag parameter

	size, err := extractFileSize(minSizeParam) // Extract size and unit type

	if err != nil {
		log.Printf("Invalid min size parameter: %v\n", err)
		return false // Invalid min size parameter
	}

	return file.Size >= size // Check if the file size is within the limit
}

func fileMimeTypeRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(*multipart.FileHeader)
	if !ok || file == nil {
		return false // Not a valid file header
	}

	mimes := strings.Split(fl.Param(), ",") // Get the allowed MIME types from the validation tag parameter
	if len(mimes) == 0 {
		return true // No MIME types specified, consider it valid
	}

	fileType := file.Header.Get("Content-Type")
	for _, mime := range mimes {
		if strings.TrimSpace(mime) == fileType {
			return true // File type matches one of the allowed MIME types
		}
	}

	return false // File type does not match any of the allowed MIME types
}
