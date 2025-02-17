package mixins

import (
	"net/http"

	"gorim.org/gorim"
)


type RetrieveMixin[T any] struct {
	GenericViewSet[T]
}

func NewRetrieveMixin[T any](
	genericViewSet GenericViewSet[T],
) *RetrieveMixin[T] {
	return &RetrieveMixin[T]{
		GenericViewSet: genericViewSet,
	}
}

func (h *RetrieveMixin[T]) Retrieve(c gorim.Context) error {
	instance := h.Child.GetObject()
	return c.JSON(http.StatusOK, instance)
}
