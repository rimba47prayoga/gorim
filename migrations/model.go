package migrations

import "time"

type GorimMigrations struct {
	ID          uint            `gorm:"primarykey"`
	Name		string			`gorm:"type:varchar(255);unique;not null"`
	Version		string			`gorm:"type:text"`
	CreatedAt   time.Time       `gorm:"type:timestamp"`
	UpdatedAt   *time.Time      `gorm:"type:timestamp"`
}

func (m GorimMigrations) TableName() string {
	return "gorim_migrations"
}
