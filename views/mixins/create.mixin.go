package mixins

import (
	"net/http"

	"gorim.org/gorim"
)


type CreateMixin[T any] struct {
	GenericViewSet[T]
}

func NewCreateMixin[T any](
	genericViewSet GenericViewSet[T],
) *CreateMixin[T] {
	return &CreateMixin[T]{
		GenericViewSet: genericViewSet,
	}
}

// @Router [POST] /api/v1/{feature}
func (h *CreateMixin[T]) Create(
	c gorim.Context,
) error {
	serializer := *h.Child.GetSerializer()
	if !serializer.IsValid() {
		return c.JSON(http.StatusBadRequest, serializer.GetErrors())
	}
	data := serializer.Create()
	return c.JSON(http.StatusCreated, data)
}
