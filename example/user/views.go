package user

import (
	"example/db"

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
