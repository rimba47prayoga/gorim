package views

import (
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/filters"
	"github.com/rimba47prayoga/gorim.git/permissions"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"github.com/rimba47prayoga/gorim.git/views/mixins"
	"gorm.io/gorm"
)

type ModelViewSetParams[T any] struct {
	QuerySet		*gorm.DB
	PKField			string
	Serializer		serializers.IModelSerializer[T]
	Filter			filters.IFilterSet
	Permissions		[]permissions.IPermission
	Child			mixins.IGenericViewSet[T]
}


type ModelViewSet[T any] struct {
	mixins.GenericViewSet[T]
	mixins.CreateMixin[T]
	Child	mixins.IGenericViewSet[T]
}

func NewModelViewSet[T any](
	params	ModelViewSetParams[T],
) *ModelViewSet[T] {
	genericViewSetParams := mixins.GenericViewSetParams[T]{
		QuerySet: params.QuerySet,
		PKField: params.PKField,
		Serializer: params.Serializer,
		Filter: params.Filter,
		Permissions: params.Permissions,
		Child: params.Child,
	}
	genericViewSet := mixins.NewGenericViewSet(genericViewSetParams)
	createMixin := mixins.NewCreateMixin[T](*genericViewSet)
	return &ModelViewSet[T]{
		GenericViewSet: *genericViewSet,
		CreateMixin: *createMixin,
	}
}

// @Router [GET] /api/v1/{feature}
func (h *ModelViewSet[T]) List(
	c gorim.Context,
) error {

	results := h.GetModelSlice()
	resultsAddr := results.Addr().Interface() //  its like &[]models.User
	queryset, err := h.FilterQuerySet(c, resultsAddr, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, gorim.Response{
			"error": err.Error(),
		})
	}
	paginate := h.PaginateQuerySet(c, queryset, resultsAddr)

	return c.JSON(http.StatusOK, paginate.GetPaginatedResponse())
}

func (h *ModelViewSet[T]) Retrieve(c gorim.Context) error {
	instance := h.GetObject()
	return c.JSON(http.StatusOK, instance)
}


// @Router [PUT] /api/v1/{feature}/:id
func (h *ModelViewSet[T]) Update(
	c gorim.Context,
) error {
	instance := h.GetObject()
	serializer := *h.GetSerializer()
	if !serializer.IsValid() {
		return c.JSON(http.StatusBadRequest, serializer.GetErrors())
	}
	data := serializer.Update(instance)
	return c.JSON(http.StatusOK, data)
}
