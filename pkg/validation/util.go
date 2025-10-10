package validation

import "encoding/json"

type InputField = string
type ErrorMessage = string

type ValidationError = map[InputField]ErrorMessage

func ParseValidationErrors(err error) (ValidationError, error) {
	bytes, err := json.Marshal(err)

	if err != nil {
		return nil, err
	}

	var validationErrMap ValidationError
	err = json.Unmarshal(bytes, &validationErrMap)

	if err != nil {
		return nil, err
	}

	return validationErrMap, nil
}
