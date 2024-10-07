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
	queryset := db.DB.Model(model)
	serializer := UserSerializer{}
	modelViewSet := views.NewModelViewSet[models.User](
		&model,
		queryset,
		&serializer,
		nil,
	)
	return &UserViewSet{
		ModelViewSet: *modelViewSet,
	}
}
