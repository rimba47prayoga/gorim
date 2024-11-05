package serializers

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/conf"
	"github.com/rimba47prayoga/gorim.git/errors"
	"github.com/rimba47prayoga/gorim.git/utils"
	"gorm.io/gorm"
)


type IModelSerializer[T any] interface {
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
	Serializer
	Child			IModelSerializer[T]  // TODO: change to interface
}

func (s *ModelSerializer[T]) Model() *T {
	var instance T
	return &instance
}

func (s *ModelSerializer[T]) Fields() []string {
	return s.GetFields()
}

func (s *ModelSerializer[T]) GetFields() []string {
	serializer := s.Child
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

func (s *ModelSerializer[T]) DB() *gorm.DB {
	return conf.DB
}

func (s *ModelSerializer[T]) SetChild(child IModelSerializer[T]) {
	s.Child = child
}

func (s *ModelSerializer[T]) SetModelAttr(model *T) {
	serializer := s.Child
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

func (s *ModelSerializer[T]) Create() *T {
	serializer := s.Child
	model := serializer.Model()
	s.SetModelAttr(model)
	serializer.DB().Create(model)
	return model
}

func (s *ModelSerializer[T]) Update(instance *T) *T {
	serializer := s.Child
	s.SetModelAttr(instance)
	serializer.DB().Save(instance)
	return instance
}
