package user

import (
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/serializers"
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
		Child: &viewset,
	}
	modelViewSet := views.NewModelViewSet(params)
	viewset.ModelViewSet = *modelViewSet
	return &viewset
}

func (h *UserViewSet) GetSerializerStruct() serializers.IModelSerializer[models.User] {
	if h.Action == "Profile" {
		return &UserProfileSerializer{}
	}
	return h.GenericViewSet.GetSerializerStruct()
}

func (h *UserViewSet) Profile(ctx gorim.Context) error {
	serializer := *h.GetSerializer()
	if !serializer.IsValid() {
		return ctx.JSON(http.StatusBadRequest, gorim.Response{
			"error": serializer.GetErrors(),
		})
	}
	data := serializer.Create()
	return ctx.JSON(http.StatusOK, data)
	
}
