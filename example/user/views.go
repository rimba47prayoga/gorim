package user

import (
	"fmt"
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/permissions"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"github.com/rimba47prayoga/gorim.git/views"
	"github.com/rimba47prayoga/gorim.git/views/mixins"
)


type ProductViewSet struct {
	mixins.GenericViewSet[models.User]
	mixins.ListMixin[models.User]
}

func NewProductViewSet() *ProductViewSet {
	viewset := ProductViewSet{}
	params := mixins.GenericViewSetParams[models.User]{
		Child: &viewset,
		Permissions: []interfaces.IPermission{
			&permissions.AllowAny{},
		},
	}
	genericViewSet := mixins.NewGenericViewSet[models.User](params)
	listMixin := mixins.NewListMixin[models.User](*genericViewSet)
	viewset.GenericViewSet = *genericViewSet
	viewset.ListMixin = *listMixin
	return &viewset
}


type UserViewSet struct {
	views.ModelViewSet[models.User]
}

// TODO: How if method List can hit without auth
func NewUserViewSet() *UserViewSet {
	var serializer UserSerializer
	viewset := UserViewSet{}
	params := views.ModelViewSetParams[models.User]{
		Serializer: &serializer,
		Child: &viewset,
	}
	modelViewSet := views.NewModelViewSet(params)
	viewset.ModelViewSet = *modelViewSet
	return &viewset
}

// override GetPermissions
func (h *UserViewSet) GetPermissions(c gorim.Context) []interfaces.IPermission {
	if h.Action == "List" {
		return []interfaces.IPermission{
			&permissions.AllowAny{},
		}
	}
	return h.GenericViewSet.GetPermissions(c)
}

// func (h *UserViewSet) GetObject() *models.User {
// 	if h.Action == "Profile" {
// 		queryset := h.GetQuerySet().Preload("Profile")
// 		return utils.GetObjectOr404[models.User](queryset, "id = ?", h.Context.Param("pk"))
// 	}
// 	return h.GenericViewSet.GetObject()
// }

func (h *UserViewSet) List(ctx gorim.Context) error {
	fmt.Println(h.GetPermissions(ctx))
	return h.ListMixin.List(ctx)
}

func (h *UserViewSet) GetSerializerStruct() serializers.IModelSerializer[models.User] {
	if h.Action == "Profile" {
		return &UserProfileSerializer{}
	}
	return h.GenericViewSet.GetSerializerStruct()
}

func (h *UserViewSet) Profile(ctx gorim.Context) error {
	instance := h.GetObject()
	serializer := *h.GetSerializer()
	if !serializer.IsValid() {
		return ctx.JSON(http.StatusBadRequest, gorim.Response{
			"error": serializer.GetErrors(),
		})
	}
	data := serializer.Update(instance)
	return ctx.JSON(http.StatusOK, data)
}
