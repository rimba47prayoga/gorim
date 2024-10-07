package views

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/filters"
	"github.com/rimba47prayoga/gorim.git/pagination"
	"github.com/rimba47prayoga/gorim.git/permissions"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"gorm.io/gorm"
)

type ActionType func(echo.Context) error


type ModelViewSet[T any] struct {
	Model			*T
	QuerySet		*gorm.DB
	Serializer		serializers.IModelSerializer[T]
	Filter			filters.IFilterSet
	Permissions		[]permissions.IPermission
	Action			string
	Context			echo.Context
	ExtraActions	[]ActionType
}

func NewModelViewSet[T any](
	model *T,
	querySet *gorm.DB,
	serializer serializers.IModelSerializer[T],
	filter	filters.IFilterSet,
) *ModelViewSet[T] {
	return &ModelViewSet[T]{
		Model: model,
		QuerySet: querySet,
		Serializer: serializer,
		Filter: filter,
	}
}

func (h *ModelViewSet[T]) RegisterAction(method ActionType) {
	h.ExtraActions = append(h.ExtraActions, method)
}

func (h *ModelViewSet[T]) GetPermissions(c echo.Context) []permissions.IPermission {
	return h.Permissions
}

func (h *ModelViewSet[T]) HasPermission(c echo.Context) bool {
	permissions := h.GetPermissions(c)
	for _, permission := range permissions {
		if !permission.HasPermission(c) {
			return false
		}
	}
	return true
}


func (h *ModelViewSet[T]) SetContext(c echo.Context) {
	h.Context = c
}


func (h *ModelViewSet[T]) SetAction(name string) {
	h.Action = name
}

func(h *ModelViewSet[T]) GetSerializer(c *echo.Context) *serializers.IModelSerializer[T] {
	return &h.Serializer
}

func (h *ModelViewSet[T]) GetQuerySet(c echo.Context) *gorm.DB {
	if h.Action == "ListDeleted" {
		return h.QuerySet.Unscoped().Where("deleted_at IS NOT NULL")
	}
	return h.QuerySet
}

func (h *ModelViewSet[T]) GetModelSlice() reflect.Value {
	// it will dynamically return slice of model specified in BaseHandler.Model
	// example: []models.User
	// Create a slice of the model type dynamically
	typeOf := reflect.TypeOf(h.Model)
	sliceType := reflect.SliceOf(typeOf)
	results := reflect.New(sliceType).Elem()
	return results
}

func (h *ModelViewSet[T]) FilterQuerySet(
	c echo.Context,
	results interface{},
	queryset *gorm.DB,
) (*gorm.DB, error) {
	if queryset == nil {
		queryset = h.GetQuerySet(c)
	}

	if h.Filter == nil {
		return queryset, nil
	}
	if err := c.Bind(h.Filter); err != nil {
		c.JSON(http.StatusBadRequest, gorim.Response{"error": err.Error()})
		return nil, err
	}
	queryset = h.Filter.ApplyFilters(h.Filter, c, queryset)

	err := queryset.Find(results).Error
	if err != nil {
		return nil, err
	}
	return queryset, nil
}

func (h *ModelViewSet[T]) PaginateQuerySet(
	ctx echo.Context,
	queryset *gorm.DB,
	results interface{},
) *pagination.Pagination {
	pagination := pagination.InitPagination(ctx, queryset)
	pagination.PaginateQuery(results)
	return pagination
}

// @Router [GET] /api/v1/{feature}
func (h *ModelViewSet[T]) List(
	c echo.Context,
) error {

	results := h.GetModelSlice()
	resultsAddr := results.Addr().Interface() //  its like &[]models.User
	queryset, err := h.FilterQuerySet(c, resultsAddr, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gorim.Response{
			"error": err.Error(),
		})
	}
	fmt.Println(queryset, resultsAddr)
	paginate := h.PaginateQuerySet(c, queryset, resultsAddr)

	return c.JSON(http.StatusOK, paginate.GetPaginatedResponse())
}
