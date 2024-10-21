package serializers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/errors"
	"github.com/rimba47prayoga/gorim.git/utils"
	"gorm.io/gorm"
)


type Meta[T any] struct {
	Model	T
	Fields	[]string
	DB		*gorm.DB
}


type IModelSerializer[T any] interface {
	IsValid() bool
	GetErrors() []errors.ValidationError
	GetContext() echo.Context
	SetContext(echo.Context)
	Meta() *Meta[T]
	SetMeta(*Meta[T])
	GetMeta() *Meta[T]
	SetChild(IModelSerializer[T])
	Create() *T
	Update(*T) *T
}


type ModelSerializer[T any] struct {
	Serializer
	Child			IModelSerializer[T]  // TODO: change to interface
	meta			*Meta[T]	`json:"-"`
}

func (s *ModelSerializer[T]) Meta() *Meta[T] {
	panic("NotImplementedError: Meta() must be overriden.")
}

func (s *ModelSerializer[T]) SetMeta(meta *Meta[T]) {
	s.meta = meta
}

func (s *ModelSerializer[T]) GetMeta() *Meta[T] {
	return s.meta
}

func (s *ModelSerializer[T]) SetChild(child IModelSerializer[T]) {
	s.Child = child
}

func (s *ModelSerializer[T]) SetModelAttr(model *T) error {
	serializer := s.Child
	for _, field := range serializer.GetMeta().Fields {

		value, err := utils.GetStructValue(serializer, field)
		fmt.Println(field, value)
		if err != nil {
			return err
		}
		utils.SetStructValue(model, field, value)
	}
	return nil
}

func (s *ModelSerializer[T]) Create() *T {
	serializer := s.Child
	meta := serializer.GetMeta()
	err := s.SetModelAttr(&meta.Model)
	if err != nil {
		panic(&errors.InternalServerError{
			Message: err.Error(),
		})
	}
	meta.DB.Create(&meta.Model)
	return &meta.Model
}

func (s *ModelSerializer[T]) Update(instance *T) *T {
	serializer := s.Child
	meta := serializer.GetMeta()
	err := s.SetModelAttr(instance)
	if err != nil {
		panic(&errors.InternalServerError{
			Message: err.Error(),
		})
	}
	meta.DB.Save(instance)
	return instance
}
