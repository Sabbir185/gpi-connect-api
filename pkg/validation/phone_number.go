package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func ValidateGlobalPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}
