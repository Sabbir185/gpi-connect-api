package validation

import "github.com/go-playground/validator/v10"

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("phone", ValidateGlobalPhone)
	return v
}

func Struct(payload any) error {
	return validate.Struct(payload)
}
