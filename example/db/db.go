package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rimba47prayoga/gorim.git/models"
	"github.com/rimba47prayoga/gorim.git/settings"
)

var DB *gorm.DB

type PostgreConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func NewPostgreConfig() *PostgreConfig {
	return &PostgreConfig{
		Host:     settings.DATABASE_POSTGRESQL_HOST,
		Port:     settings.DATABASE_POSTGRESQL_PORT,
		User:     settings.DATABASE_POSTGRESQL_USER,
		Password: settings.DATABASE_POSTGRESQL_PASSWORD,
		DbName:   settings.DATABASE_POSTGRESQL_DB_NAME,
	}
}

func SetupDatabase() {
	config := NewPostgreConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		panic(err)
	}

	DB = db
}
