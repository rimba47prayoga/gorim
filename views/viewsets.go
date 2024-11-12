package views

import (
	"github.com/rimba47prayoga/gorim.git/filters"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"github.com/rimba47prayoga/gorim.git/views/mixins"
	"gorm.io/gorm"
)

type ModelViewSetParams[T any] struct {
	QuerySet		*gorm.DB
	PKField			string
	Serializer		serializers.IModelSerializer[T]
	Filter			filters.IFilterSet
	Permissions		[]interfaces.IPermission
	Child			mixins.IGenericViewSet[T]
}


type ModelViewSet[T any] struct {
	mixins.GenericViewSet[T]
	mixins.CreateMixin[T]
	mixins.RetrieveMixin[T]
	mixins.UpdateMixin[T]
	mixins.ListMixin[T]
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
	updateMixin := mixins.NewUpdateMixin[T](*genericViewSet)
	retrieveMixin := mixins.NewRetrieveMixin[T](*genericViewSet)
	listMixin := mixins.NewListMixin[T](*genericViewSet)
	return &ModelViewSet[T]{
		GenericViewSet: *genericViewSet,
		CreateMixin: *createMixin,
		UpdateMixin: *updateMixin,
		RetrieveMixin: *retrieveMixin,
		ListMixin: *listMixin,
	}
}
