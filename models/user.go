package models

import (
	"github.com/rimba47prayoga/gorim.git/errors"
	"github.com/rimba47prayoga/gorim.git/utils"
)


type User struct {
	BaseModel
	// Username	string			`gorm:"type:varchar(255)" json:"username"`
	Email		string			`gorm:"type:varchar(255)" json:"email"`
	Password	string			`gorm:"type:varchar(255)" json:"password"`
}

type AbstractUser struct {
	BaseModel
	Email		string			`gorm:"type:varchar(255)" json:"email"`
	Password	string			`gorm:"type:varchar(255)" json:"password"`
}

func (m *AbstractUser) SetPassword(passwd string) {
	hashedPassword, err := utils.HashPassword(passwd)
	if err != nil {
		errors.Raise(&errors.InternalServerError{
			Message: err.Error(),
		})
	}
	m.Password = hashedPassword
}

func (m *AbstractUser) CheckPassword(passwd string) bool {
	return utils.VerifyPassword(passwd, m.Password)
}
