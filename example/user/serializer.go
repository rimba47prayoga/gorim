package user

import (
	"gorim.org/gorim/models"
	"gorim.org/gorim/serializers"
)

type UserSerializer struct {
	serializers.ModelSerializer[models.User]
	Email		string		`validate:"required,email" json:"email"`
	Password	string		`validate:"required" json:"password"`
}

type UserProfileSerializer struct {
	serializers.ModelSerializer[models.User]
	Username	string		`validate:"required" json:"username"`
	Password	string		`validate:"required" json:"password"`
}

func (s *UserProfileSerializer) ValidateUsername() {
	s.AddError("username", "not valid")
}

func (s *UserProfileSerializer) Validate() {
	s.ModelSerializer.Validate()
}
