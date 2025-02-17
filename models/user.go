package models

import (
	"gorim.org/gorim/errors"
	"gorim.org/gorim/utils"
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
