package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

var ErrInvalidPhone = errors.New("invalid phone number")

func ValidateGlobalPhone(fl validator.FieldLevel) bool {
	_, err := NormalizePhone(fl.Field().String())
	return err == nil
}

// NormalizePhone returns the E.164 form, so one number always maps to one account.
func NormalizePhone(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return "", ErrInvalidPhone
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
