package user

import (
	"fmt"
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/permissions"
	"github.com/rimba47prayoga/gorim.git/serializers"
	"github.com/rimba47prayoga/gorim.git/utils"
	"github.com/rimba47prayoga/gorim.git/views"
)


type UserViewSet struct {
	views.ModelViewSet[models.User]
}

// TODO: How if method List can hit without auth
func NewUserViewSet() *UserViewSet {
	var serializer UserSerializer
	viewset := UserViewSet{}
	params := views.ModelViewSetParams[models.User]{
		Serializer: &serializer,
		PKField: "uuid",
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

func (h *UserViewSet) GetObject() *models.User {
	if h.Action == "UpdateProfile" {
		queryset := h.GetQuerySet().Preload("Profile")
		return utils.GetObjectOr404[models.User](queryset, "id = ?", h.Context.Param("pk"))
	}
	return h.GenericViewSet.GetObject()
}

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

func (h *UserViewSet) UpdateProfile(ctx gorim.Context) error {
	profile := h.GetObject()

	serializer := *h.GetSerializer()
	if !serializer.IsValid() {
		return ctx.JSON(http.StatusBadRequest, gorim.Response{
			"error": serializer.GetErrors(),
		})
	}
	data := serializer.Update(profile)
	return ctx.JSON(http.StatusOK, data)
}
