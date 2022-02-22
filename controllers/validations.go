package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

//ValidationError -> represent error
type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//MakeValidationError -> mapping validator to json
func MakeValidationError(errorCode string, errorMessage string) []ValidationError {
	return []ValidationError{{errorCode, errorMessage}}
}

//MakeValidationErrors -> mapping validator to json
func MakeValidationErrors(errs validator.ValidationErrors) []ValidationError {
	validationErrors := make([]ValidationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, ValidationError{
			Code:    ErrorInvalidField,
			Message: fmt.Sprintf("[%s] type %s, validations errors: %s", validationErr.Namespace(), validationErr.Type().String(), validationErr.ActualTag()),
		})
	}

	return validationErrors
}
