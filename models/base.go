package models

import (
	"time"

	"gorm.io/gorm"
)


type BaseModel struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	CreatorId   *string         `gorm:"type:varchar(255)" json:"-"`
	CreatorName *string         `gorm:"type:varchar(255)" json:"-"`
	UpdaterId   *string         `gorm:"type:varchar(255)" json:"-"`
	UpdaterName *string         `gorm:"type:varchar(255)" json:"-"`
	DeleterId   *string         `gorm:"type:varchar(255)" json:"-"`
	DeleterName *string         `gorm:"type:varchar(255)" json:"-"`
	CreatedAt   time.Time       `gorm:"type:timestamp" json:"-"`
	UpdatedAt   *time.Time      `gorm:"type:timestamp" json:"-"`
	DeletedAt   *gorm.DeletedAt `gorm:"index" json:"-"`
}
