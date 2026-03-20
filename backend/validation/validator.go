package validation

import (
	"fmt"
	"smartech/backend/errors"

	"github.com/go-playground/validator/v10"
)

// Validator es el validador global
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct valida una estructura y retorna errores de validación
func ValidateStruct(data interface{}) errors.ValidationErrors {
	errs := errors.ValidationErrors{}

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs[err.Field()] = getErrorMessage(err)
		}
	}

	return errs
}

// ValidateVar valida una variable individual
func ValidateVar(value interface{}, tag string) error {
	return validate.Var(value, tag)
}

// getErrorMessage retorna un mensaje amigable para el error de validación
func getErrorMessage(err validator.FieldError) string {
	switch err.ActualTag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
	case "numeric":
		return fmt.Sprintf("%s must be numeric", err.Field())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
	case "unique":
		return fmt.Sprintf("%s must be unique", err.Field())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// IsValid retorna true si no hay errores de validación
func IsValid(data interface{}) bool {
	return len(ValidateStruct(data)) == 0
}
