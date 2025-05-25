package validation

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

type errorInputMap = map[string]string

func Validate[T any](request T) ([]errorInputMap, error) {
	validate := validator.New()

	translator, err := prepareTranslator(validate)

	if err != nil {
		return nil, err
	}

	err = validate.Struct(request)

	if err != nil {
		return mappingErrors(err, translator), nil
	}

	return nil, nil
}

func prepareTranslator(validate *validator.Validate) (ut.Translator, error) {
	english := en.New()
	uni := ut.New(english, english)

	translator, found := uni.GetTranslator("en")

	if !found {
		return nil, errors.New("translator not found")
	}

	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, errors.New("failed to register translations: " + err.Error())
	}

	return translator, nil
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
