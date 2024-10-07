package user

import (
	"example/db"

	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/serializers"
)

type UserSerializer struct {
	serializers.IModelSerializer[models.User]
}

func (s *UserSerializer) Meta() *serializers.Meta[models.User] {
	return &serializers.Meta[models.User]{
		Model: models.User{},
		Fields: []string{"__all__"},
		DB: db.DB,
	}
}
