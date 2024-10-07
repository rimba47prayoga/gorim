package serializers

import (
	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git/errors"
	"gorm.io/gorm"
)


type Meta[T any] struct {
	Model	T
	Fields	[]string
	DB		*gorm.DB
}


type IModelSerializer[T any] interface {
	Validate() bool
	IsValid() bool
	GetErrors() []errors.ValidationError
	GetContext() echo.Context
	Meta() *Meta[T]
	Save() *T
	Create() *T
	Update() *T
}


type ModelSerializer[T any] struct {
	Serializer
	Instance	*T
}

func (s *ModelSerializer[T]) Meta() *Meta[T] {
	panic("NotImplementedError: Meta() must be overriden.")
}

func (s *ModelSerializer[T]) Save() *T {
	if s.Instance == nil {
		return s.Create()
	}
	return s.Update()
}

func (s *ModelSerializer[T]) Create() *T {
	meta := s.Meta()
	meta.DB.Model(&meta.Model).Create(s)
	return &meta.Model
}

func (s *ModelSerializer[T]) Update() *T {
	meta := s.Meta()
	meta.DB.Model(&s.Instance).Save(s)
	return &meta.Model
}
