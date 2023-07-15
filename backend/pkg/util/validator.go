package util

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var Validate *validator.Validate

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define fields map.
	fields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
