package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	validate := validator.New()
	return &Validator{
		validate: validate,
	}
}

func (v *Validator) Struct(s any) error {
	return v.validate.Struct(s)
}

// ValidateStruct validates a struct and returns formatted errors
func (v *Validator) ValidateStruct(s any) map[string]any {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]any)
	for _, err := range err.(validator.ValidationErrors) {
		errors[jsonFieldName(err.Field())] = formatError(err)
	}

	return errors
}

// DecodeAndValidate decodes JSON from request and validates it
func (v *Validator) DecodeAndValidate(r *http.Request, dst any) map[string]any {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return map[string]any{
			"_error": fmt.Sprintf("Invalid JSON: %s", err.Error()),
		}
	}

	return v.ValidateStruct(dst)
}

// formatError formats a validation error into a user-friendly message
func formatError(err validator.FieldError) string {
	field := jsonFieldName(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// jsonFieldName converts struct field name to JSON field name (camelCase to snake_case)
func jsonFieldName(field string) string {
	// Simple conversion: convert CamelCase to snake_case
	var result strings.Builder
	for i, r := range field {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
