package validation

import "errors"

type Rule = string
type Validator = func(value interface{}) error

var (
	databaseRepository   DatabaseValidationRepository
	ErrValidatorNotFound = errors.New("validator not found")
	validators           = make(map[Rule]Validator)
)

func Init(db DatabaseValidationRepository) {
	databaseRepository = db
}

func GetValidator(name string) (Validator, error) {
	if v, ok := validators[name]; ok {
		return v, nil
	}
	return nil, ErrValidatorNotFound
}

func RegisterValidator(name string, validator Validator) {
	validators[name] = validator
}
