package mixins

import (
	"net/http"
	"reflect"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/errors"
	"github.com/rimba47prayoga/gorim.git/filters"
	"github.com/rimba47prayoga/gorim.git/pagination"
	"github.com/rimba47prayoga/gorim.git/permissions"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"github.com/rimba47prayoga/gorim.git/utils"
	"gorm.io/gorm"
)

type ActionType func(gorim.Context) error

type IGenericViewSet[T any] interface {
	GetSerializer() *serializers.IModelSerializer[T]
}


type GenericViewSetParams[T any] struct {
	QuerySet		*gorm.DB
	PKField			string
	Serializer		serializers.IModelSerializer[T]
	Filter			filters.IFilterSet
	Permissions		[]permissions.IPermission
	Child			IGenericViewSet[T]
}


type GenericViewSet[T any] struct {
	Model			*T
	QuerySet		*gorm.DB
	PKField			string
	Serializer		serializers.IModelSerializer[T]
	Filter			filters.IFilterSet
	Permissions		[]permissions.IPermission
	Action			string
	Context			gorim.Context
	ExtraActions	[]ActionType
	Child			IGenericViewSet[T]
}

func NewGenericViewSet[T any](
	params GenericViewSetParams[T],
) *GenericViewSet[T] {
	model := params.QuerySet.Statement.Model.(*T)
	return &GenericViewSet[T]{
		Model: model,
		QuerySet: params.QuerySet,
		Serializer: params.Serializer,
		Filter: params.Filter,
		Permissions: params.Permissions,
		Child: params.Child,
	}
}

func (h *GenericViewSet[T]) GetPKField() string {
	if h.PKField == "" {
		return "id"
	}
	return h.PKField
}

func (h *GenericViewSet[T]) GetPermissions(c gorim.Context) []permissions.IPermission {
	return h.Permissions
}

func (h *GenericViewSet[T]) HasPermission(c gorim.Context) bool {
	permissions := h.GetPermissions(c)
	for _, permission := range permissions {
		if !permission.HasPermission(c) {
			return false
		}
	}
	return true
}


func (h *GenericViewSet[T]) SetContext(c gorim.Context) {
	h.Context = c
}


func (h *GenericViewSet[T]) SetAction(name string) {
	h.Action = name
}

func(h *GenericViewSet[T]) SetupSerializer(
	serializer serializers.IModelSerializer[T],
) *serializers.IModelSerializer[T] {
	serializer.SetContext(h.Context)
	serializer.SetMeta(serializer.Meta())
	if err := h.Context.Bind(&serializer); err != nil {
		panic(&errors.InternalServerError{
			Message: err.Error(),
		})
	}
	serializer.SetChild(serializer)
	return &serializer
} 

func(h *GenericViewSet[T]) GetSerializer() *serializers.IModelSerializer[T] {
	serializer := h.GetSerializerStruct()
	return h.SetupSerializer(serializer)
}

func(h *GenericViewSet[T]) GetSerializerStruct() serializers.IModelSerializer[T] {
	return h.Serializer
}

func (h *GenericViewSet[T]) GetQuerySet() *gorm.DB {
	if h.Action == "ListDeleted" {
		return h.QuerySet.Unscoped().Where("deleted_at IS NOT NULL")
	}
	return h.QuerySet
}

func (h *GenericViewSet[T]) GetObject() *T {
	pkField := h.GetPKField()
	id := h.Context.Param("id")
	queryset := h.GetQuerySet()
	result := utils.GetObjectOr404[T](queryset, pkField + " = ?", id)
	return result
}

func (h *GenericViewSet[T]) GetModelSlice() reflect.Value {
	// it will dynamically return slice of model specified in BaseHandler.Model
	// example: []models.User
	// Create a slice of the model type dynamically
	typeOf := reflect.TypeOf(h.Model)
	sliceType := reflect.SliceOf(typeOf)
	results := reflect.New(sliceType).Elem()
	return results
}

func (h *GenericViewSet[T]) FilterQuerySet(
	c gorim.Context,
	results interface{},
	queryset *gorm.DB,
) (*gorm.DB, error) {
	if queryset == nil {
		queryset = h.GetQuerySet()
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

func (h *GenericViewSet[T]) PaginateQuerySet(
	ctx gorim.Context,
	queryset *gorm.DB,
	results interface{},
) *pagination.Pagination {
	pagination := pagination.InitPagination(ctx, queryset)
	pagination.PaginateQuery(results)
	return pagination
}
