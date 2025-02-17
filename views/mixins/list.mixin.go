package mixins

import (
	"net/http"

	"gorim.org/gorim"
)


type ListMixin[T any] struct {
	GenericViewSet[T]
}

func NewListMixin[T any](
	genericViewSet GenericViewSet[T],
) *ListMixin[T] {
	return &ListMixin[T]{
		GenericViewSet: genericViewSet,
	}
}

// @Router [GET] /api/v1/{feature}
func (h *ListMixin[T]) List(
	c gorim.Context,
) error {
	viewset := h.Child
	results := viewset.GetModelSlice()
	resultsAddr := results.Addr().Interface() //  its like &[]models.User
	queryset := viewset.FilterQuerySet(resultsAddr, nil)
	paginate := viewset.PaginateQuerySet(resultsAddr, queryset)
	return c.JSON(http.StatusOK, paginate.GetPaginatedResponse())
}
