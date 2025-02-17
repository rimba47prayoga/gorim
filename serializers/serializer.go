package serializers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorim.org/gorim/errors"
)

type ISerializer interface {
	IsValid() bool
	GetErrors() []errors.ValidationError
	GetContext() echo.Context
	SetContext(echo.Context)
}

// Serializer struct with embedded error handling
type Serializer struct {
	errors			[]errors.ValidationError
	structType		reflect.Type
	context			echo.Context
}

func (s *Serializer) GetContext() echo.Context {
	return s.context
}

func (s *Serializer) SetContext(c echo.Context) {
	s.context = c
}

func (s *Serializer) GetErrors() []errors.ValidationError {
	return s.errors
}

func (s *Serializer) AddError(field string, message string) {
	s.errors = append(s.errors, errors.ValidationError{
		Field: field,
		Message: message,
	})
}

// Function to extract the JSON or form tag name from the struct field.
func (s *Serializer) getFieldName(structType reflect.Type, fieldName string) string {
	if field, ok := structType.FieldByName(fieldName); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			return strings.Split(jsonTag, ",")[0] // Handle cases like `json:"field_name,omitempty"`
		}

		formTag := field.Tag.Get("form")
		if formTag != "" && formTag != "-" {
			return strings.Split(formTag, ",")[0]
		}
	}
	return fieldName // Default to field name if no tag is found
}

// HandleError processes and formats validation errors.
func (s *Serializer) HandleError(err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var validationErrors []errors.ValidationError
		for _, e := range errs {
			fieldName := s.getFieldName(s.structType, e.StructField())
			validationErrors = append(validationErrors, errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s is %s", fieldName, e.Tag()),
			})
		}
		s.errors = validationErrors
	} else {
		s.errors = append(s.errors, errors.ValidationError{
			Field:   reflect.TypeOf(err).String(),
			Message: err.Error(),
		})
	}
}


// IsValid validates the serializer and handles errors.
func (s *Serializer) IsValid() bool {
	s.structType = reflect.TypeOf(s).Elem()
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		s.HandleError(err)
		return false
	}
	return true
}
