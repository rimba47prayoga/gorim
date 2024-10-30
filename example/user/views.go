package user

import (
	"example/db"
	"net/http"

	"github.com/rimba47prayoga/gorim.git"
	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/views"
)


type UserViewSet struct {
	views.ModelViewSet[models.User]
}

func NewUserViewSet() *UserViewSet {
	var model models.User
	var serializer UserSerializer
	queryset := db.DB.Model(&model)
	viewset := UserViewSet{}
	params := views.ModelViewSetParams[models.User]{
		QuerySet: queryset,
		Serializer: &serializer,
		Child: &viewset,
	}
	modelViewSet := views.NewModelViewSet(params)
	viewset.ModelViewSet = *modelViewSet
	return &viewset
}

func (h *UserViewSet) Profile(ctx gorim.Context) error {
	return ctx.JSON(http.StatusOK, gorim.Response{
		"status": true,
	})
}
