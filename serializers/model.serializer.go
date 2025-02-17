package serializers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorim.org/gorim/conf"
	"gorim.org/gorim/errors"
	"gorim.org/gorim/utils"
	"gorm.io/gorm"
)


type IModelSerializer[T any] interface {
	Validate()
	IsValid() bool
	GetErrors() []errors.ValidationError
	GetContext() echo.Context
	SetContext(echo.Context)
	SetChild(IModelSerializer[T])
	Model() *T
	Fields() []string
	DB() *gorm.DB
	Create() *T
	Update(*T) *T
}


type ModelSerializer[T any] struct {
	errors			[]errors.ValidationError
	context			echo.Context
	child			IModelSerializer[T]
}

// ------ Metadata ------
func (s *ModelSerializer[T]) Model() *T {
	var instance T
	return &instance
}

func (s *ModelSerializer[T]) Fields() []string {
	return s.GetFields()
}

func (s *ModelSerializer[T]) DB() *gorm.DB {
	return conf.DB
}
// ------ END ------

// ------ Getters ------
func (s *ModelSerializer[T]) GetContext() echo.Context {
	return s.context
}

func (s *ModelSerializer[T]) GetErrors() []errors.ValidationError {
	return s.errors
}

func (s *ModelSerializer[T]) GetFields() []string {
	serializer := s.child
	val := reflect.ValueOf(serializer)
	fields := []string{}
	typ := val.Type()

	// Ensure the input is a struct
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fields // Return empty slice if not a struct
	}

	// Iterate over fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			fields = append(fields, field.Name)
		}
	}

	return fields
}

// Function to extract the JSON or form tag name from the struct field.
func (s *ModelSerializer[T]) GetFieldName(fieldName string) string {
	structType := reflect.TypeOf(s.child).Elem()
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
// ------ END ------

// ------ Setters ------
func (s *ModelSerializer[T]) SetContext(c echo.Context) {
	s.context = c
}

func (s *ModelSerializer[T]) SetChild(child IModelSerializer[T]) {
	s.child = child
}

func (s *ModelSerializer[T]) SetModelAttr(model *T) {
	serializer := s.child
	for _, field := range serializer.Fields() {

		value, err := utils.GetStructValue(serializer, field)
		if err != nil {
			errors.Raise(&errors.InternalServerError{
				Message: err.Error(),
			})
		}
		err = utils.SetStructValue(model, field, value)
		if err != nil {
			errors.Raise(&errors.InternalServerError{
				Message: err.Error(),
			})
		}
	}
}
// ------ END ------

// ------ Error Handlers ------
func (s *ModelSerializer[T]) AddError(field string, message string) {
	s.errors = append(s.errors, errors.ValidationError{
		Field: field,
		Message: message,
	})
}

// HandleError processes and formats validation errors.
func (s *ModelSerializer[T]) HandleError(err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var validationErrors []errors.ValidationError
		for _, e := range errs {
			fieldName := s.GetFieldName(e.StructField())
			validationErrors = append(validationErrors, errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("%s is %s", fieldName, e.Tag()),
			})
		}
		s.errors = validationErrors
	} else {
		s.errors = append(s.errors, errors.ValidationError{
			Field:   "non_field_errors",
			Message: err.Error(),
		})
	}
}
// ------ END ------

// ------ Validation ------
func (s *ModelSerializer[T]) Validate() {
	serializer := s.child
	validate := validator.New()
	if err := validate.Struct(serializer); err != nil {
		s.HandleError(err)
		return
	}
	s.ValidateField()
}

func (s *ModelSerializer[T]) ValidateField() {
	serializer := s.child
	serializerVal := reflect.ValueOf(serializer)
	for _, field := range serializer.Fields() {
		methodName := fmt.Sprintf("Validate%s", field)
		if utils.HasAttr(serializer, methodName) {
			methodVal := serializerVal.MethodByName(methodName)
			methodVal.Call([]reflect.Value{})
		}
	}
}

// IsValid validates the serializer and handles errors.
func (s *ModelSerializer[T]) IsValid() bool {
	s.child.Validate()
	isValid := len(s.errors) == 0
	return isValid
}
// ------ END ------

func (s *ModelSerializer[T]) Create() *T {
	serializer := s.child
	model := serializer.Model()
	s.SetModelAttr(model)
	serializer.DB().Create(model)
	return model
}

func (s *ModelSerializer[T]) Update(instance *T) *T {
	serializer := s.child
	s.SetModelAttr(instance)
	serializer.DB().Save(instance)
	return instance
}
