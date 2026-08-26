package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	errMap := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			field := fieldErr.Field()
			param := fieldErr.Param()

			switch fieldErr.Tag() {
			case "required":
				errMap[field] = fmt.Sprintf("%s is required", field)
			case "required_if":
				errMap[field] = fmt.Sprintf("%s is required when %s matches", field, param)
			case "required_unless":
				errMap[field] = fmt.Sprintf("%s is required unless %s is provided", field, param)
			case "min":
				errMap[field] = fmt.Sprintf("%s must be at least %s characters", field, param)
			case "max":
				errMap[field] = fmt.Sprintf("%s must not exceed %s characters", field, param)
			case "len":
				errMap[field] = fmt.Sprintf("%s must be exactly %s characters long", field, param)
			case "email":
				errMap[field] = fmt.Sprintf("%s must be a valid email address", field)
			case "url":
				errMap[field] = fmt.Sprintf("%s must be a valid URL", field)
			case "uri":
				errMap[field] = fmt.Sprintf("%s must be a valid URI", field)
			case "uuid", "uuid4":
				errMap[field] = fmt.Sprintf("%s must be a valid UUID", field)
			case "numeric", "number":
				errMap[field] = fmt.Sprintf("%s must be a valid number", field)
			case "alpha":
				errMap[field] = fmt.Sprintf("%s must contain only letters", field)
			case "alphanumeric":
				errMap[field] = fmt.Sprintf("%s must contain only letters and numbers", field)
			case "oneof":
				errMap[field] = fmt.Sprintf("%s must be one of: [%s]", field, param)
			case "gt":
				errMap[field] = fmt.Sprintf("%s must be greater than %s", field, param)
			case "gte":
				errMap[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
			case "lt":
				errMap[field] = fmt.Sprintf("%s must be less than %s", field, param)
			case "lte":
				errMap[field] = fmt.Sprintf("%s must be less than or equal to %s", field, param)
			case "eqfield":
				errMap[field] = fmt.Sprintf("%s must match %s", field, param)
			case "nefield":
				errMap[field] = fmt.Sprintf("%s must not match %s", field, param)
			case "datetime":
				errMap[field] = fmt.Sprintf("%s must match date format %s", field, param)
			case "e164":
				errMap[field] = fmt.Sprintf("%s must be a valid phone number with country code", field)
			case "ip", "ipv4", "ipv6":
				errMap[field] = fmt.Sprintf("%s must be a valid IP address", field)
			case "boolean":
				errMap[field] = fmt.Sprintf("%s must be a boolean value", field)
			default:
				errMap[field] = fmt.Sprintf("%s failed on '%s' validation", field, fieldErr.Tag())
			}
		}
	}
	return errMap
}
