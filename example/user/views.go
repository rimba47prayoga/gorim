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
	queryset := db.DB.Model(model)
	modelViewSet := views.NewModelViewSet(
		&model,
		queryset,
		&serializer,
		nil,
	)
	return &UserViewSet{
		ModelViewSet: *modelViewSet,
	}
}
