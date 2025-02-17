package settings

import (
	"example/migrations"
	"fmt"
	"time"

	"gorim.org/gorim"
	"gorim.org/gorim/conf"
	"gorim.org/gorim/interfaces"
	"gorim.org/gorim/middlewares"
	"gorim.org/gorim/permissions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// its just for flag to check if settings was configured.
var CONFIGURED bool

var DEBUG bool

var ALLOWED_HOSTS []string = []string{"localhost"}

var DATABASE conf.Database

var Server *gorim.Server

var HOST string
var PORT uint

func SetupDatabase() *gorm.DB {
	config := DATABASE
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
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
	return db
}

func SetupMiddlewares() {
	Server.Use(&middlewares.LoggerMiddleware{})
	Server.Use(&middlewares.AllowedHostsMiddleware{})
	// add other middlewares here.

	// keep this at bottom.
	Server.Use(&middlewares.RecoverMiddleware{})
}

func Configure() {
	conf.UseEnv(".env")
	CONFIGURED = true
	DEBUG = true
	DATABASE = conf.Database{
		Name: "example_gorim",
		Host: "localhost",
		Port: 5432,
		User: "rimbaprayoga",
		Password: "qweqweqwe",
	}
	HOST = "127.0.0.1"
	PORT = 8000
	Server = gorim.New()
	db := SetupDatabase()
	SetupMiddlewares()

	conf.DEFAULT_PERMISSION_STRUCTS = []interfaces.IPermission{
		&permissions.AllowAny{},
	}
	conf.ALLOWED_HOSTS = ALLOWED_HOSTS

	// its for gorim settings.
	conf.DB = db
	conf.GorimServer = Server
	conf.HOST = HOST
	conf.PORT = PORT
	conf.MigrationInstance = migrations.MigrationInstance
}
