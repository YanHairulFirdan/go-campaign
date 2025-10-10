package v1

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	validationPkg "go-campaign.com/pkg/validation"
)

type UserRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserRegisterRequest) Validate() error {
	// Implement validation logic here if needed
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&r.Email, validation.Required, validation.Length(5, 100), is.Email, validationPkg.Unique("users", "email", "", nil, "Email already taken")),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
	)
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *UserLoginRequest) Validate() error {
	// Implement validation logic here if needed
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Email, validation.Required, is.Email, validation.Length(5, 100)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 100)),
	)
}
