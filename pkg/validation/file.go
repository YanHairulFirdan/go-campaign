package validation

import (
	"errors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func fileRequiredRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(multipart.FileHeader)

	return ok && file.Size > 0 // Check if the file is present and has a size greater than 0
}

func fileMaxSizeRule(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
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
	file, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
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
	file, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false // Not a valid file header
	}

	// log param
	log.Printf("param %v", fl.Param())
	mimes := strings.Split(fl.Param(), "-") // Get the allowed MIME types from the validation tag parameter
	log.Printf("Allowed MIME types: %v\n", mimes)
	if len(mimes) == 0 {
		return true // No MIME types specified, consider it valid
	}

	fileType := file.Header.Get("Content-Type")
	log.Printf("file type: %v", fileType)
	for _, mime := range mimes {
		if strings.TrimSpace(mime) == fileType {
			return true // File type matches one of the allowed MIME types
		}
	}

	return false // File type does not match any of the allowed MIME types
}
