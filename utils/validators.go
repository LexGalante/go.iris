package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

//ValidatorLicensePlate -> custom validator for licenseplate
func ValidatorLicensePlate(fieldLevel validator.FieldLevel) bool {
	match, _ := regexp.MatchString("[A-Z]{3}[0-9][0-9A-Z][0-9]{2}", fieldLevel.Field().String())

	return match
}
