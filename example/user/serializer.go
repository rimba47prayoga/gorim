package user

import (
	"example/db"

	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/serializers"
)

type UserSerializer struct {
	serializers.ModelSerializer[models.User]
	Email		string		`validate:"required,email" json:"email"`
	Password	string		`validate:"required" json:"password"`
}

func (s *UserSerializer) Meta() *serializers.Meta[models.User] {
	return &serializers.Meta[models.User]{
		Model: models.User{},
		Fields: []string{"Email", "Password"},
		DB: db.DB,
	}
}

// func (s *UserSerializer) Create() *models.User {
// 	var user models.User
// 	s.SetModelAttr(&user)
// 	db.DB.Model(&user).Save(s)
// 	return &user
// }
