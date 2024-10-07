package models


type User struct {
	BaseModel
	Email		string			`gorm:"type:varchar(255)" json:"email"`
	Password	string			`gorm:"type:varchar(255)" json:"password"`
}
