package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorMessage(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		messages := make([]string, 0, len(validationErrors))

		for _, fieldError := range validationErrors {
			field := fieldError.Field()

			switch fieldError.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", field))
			case "uuid":
				messages = append(messages, fmt.Sprintf("%s must be a valid UUID", field))
			case "email":
				messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
			case "url":
				messages = append(messages, fmt.Sprintf("%s must be a valid URL", field))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s failed %s validation", field, fieldError.Tag()))
			}
		}

		return strings.Join(messages, ", ")
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		var expected string

		switch unmarshalTypeError.Type.Kind() {
		case reflect.String:
			expected = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			expected = "integer"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			expected = "integer"
		case reflect.Float32, reflect.Float64:
			expected = "number"
		case reflect.Bool:
			expected = "boolean"
		default:
			expected = unmarshalTypeError.Type.String()
		}

		return fmt.Sprintf("%s must be a %s", unmarshalTypeError.Field, expected)
	}

	return "invalid body"
}
