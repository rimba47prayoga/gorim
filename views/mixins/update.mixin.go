package mixins

import (
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
)


type UpdateMixin[T any] struct {
	GenericViewSet[T]
}


func NewUpdateMixin[T any](
	genericViewSet GenericViewSet[T],
) *UpdateMixin[T] {
	return &UpdateMixin[T]{
		GenericViewSet: genericViewSet,
	}
}


// @Router [PUT] /api/v1/{feature}/:id
func (h *UpdateMixin[T]) Update(
	c gorim.Context,
) error {
	instance := h.Child.GetObject()
	serializer := *h.Child.GetSerializer()
	if !serializer.IsValid() {
		return c.JSON(http.StatusBadRequest, serializer.GetErrors())
	}
	data := serializer.Update(instance)
	return c.JSON(http.StatusOK, data)
}
