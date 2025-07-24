package validation

import (
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
	databaseRepository        DatabaseValidationRepository
)

func Init(db DatabaseValidationRepository) error {
	validate = validator.New()

	err := prepareTranslator()

	if err != nil {
		return err
	}

	registerCustomRules()
	databaseRepository = db

	return nil
}

func registerCustomRules() {
	validate.RegisterValidation("unique", uniqueRule)

	validate.RegisterTranslation("unique", translator, func(ut ut.Translator) error {
		return ut.Add("unique", "{0} must be unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		fieldName := fe.Field()
		return fmt.Sprintf("%s must be unique", fieldName)
	})

	validate.RegisterValidation("file_required", fileRequiredRule)
	validate.RegisterTranslation("file_required", translator, func(ut ut.Translator) error {
		return ut.Add("file_required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		fieldName := fe.Field()
		return fmt.Sprintf("%s is required", fieldName)
	})

	validate.RegisterValidation("file_max_size", fileMaxSizeRule)
	validate.RegisterTranslation("file_max_size", translator, func(ut ut.Translator) error {
		return ut.Add("file_max_size", "{0} exceeds the maximum allowed size", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		fieldName := fe.Field()
		return fmt.Sprintf("%s exceeds the maximum allowed size", fieldName)
	})

	validate.RegisterValidation("file_min_size", fileMinSizeRule)
	validate.RegisterTranslation("file_min_size", translator, func(ut ut.Translator) error {
		return ut.Add("file_min_size", "{0} is below the minimum required size", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		fieldName := fe.Field()
		return fmt.Sprintf("%s is below the minimum required size", fieldName)
	})

	validate.RegisterValidation("file_mime_type", fileMimeTypeRule)
	validate.RegisterTranslation("file_mime_type", translator, func(ut ut.Translator) error {
		return ut.Add("file_mime_type", "{0} has an invalid file type", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {

		fieldName := fe.Field()
		return fmt.Sprintf("%s has an invalid file type", fieldName)
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
	translator, found = uni.GetTranslator("en")

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
