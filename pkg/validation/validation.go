package validation

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

type errorInputMap = map[string]string

type validationException struct {
	Column string
	Value  any
}

type ValidationExceptions map[string]validationException

var (
	validate                  *validator.Validate
	translator                ut.Translator
	validationFieldExceptions ValidationExceptions
)

func Init(db *sql.DB) error {
	validate = validator.New()

	err := prepareTranslator()

	if err != nil {
		return err
	}

	registerCustomRules(db)

	return nil
}

func registerCustomRules(db *sql.DB) {
	// Here you can register custom validation rules if needed
	// For example:
	// validate.RegisterValidation("unique", UniqueRule(db))
	// This is where you would add any custom validation logic
	// such as checking for unique values in the database.

	validate.RegisterValidation("unique", uniqueRule(db))

	validate.RegisterTranslation("unique", translator, func(ut ut.Translator) error {
		return ut.Add("unique", "{0} must be unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		// Return a custom error message instead of calling fe.Translate(ut)
		fieldName := fe.Field() // Get the field name
		return fmt.Sprintf("%s must be unique", fieldName)
	})
}

func Validate[T any](request T, incomingExceptions ValidationExceptions) ([]errorInputMap, error) {
	validationFieldExceptions = incomingExceptions

	err := validate.Struct(request)

	if err != nil {
		return mappingErrors(err, translator), nil
	}

	return nil, nil
}

func prepareTranslator() error {
	english := en.New()
	uni := ut.New(english, english)

	var found bool
	translator, found = uni.GetTranslator("en") // Assign to the global translator variable

	if !found {
		return errors.New("translator not found")
	}

	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		return fmt.Errorf("failed to register translations: %w", err)
	}

	return nil
}

func mappingErrors(err error, translator ut.Translator) []errorInputMap {
	var errs validator.ValidationErrors

	errors.As(err, &errs)

	validationErrors := make([]errorInputMap, 0)

	for _, validationError := range errs {
		validationErrors = append(validationErrors, map[string]string{
			"field":   validationError.Field(),
			"message": validationError.Translate(translator),
		})
	}

	return validationErrors
}
