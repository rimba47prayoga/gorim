package conf

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rimba47prayoga/gorim.git/interfaces"
	"gorm.io/gorm"
)

var DB *gorm.DB
var GorimServer interface{}  // use type assertion

type Database struct {
	Name		string
	Host		string
	Port		int
	User		string
	Password	string
}

var ENV_PATH = ".env"
var HOST = "http://localhost:8000/"
var PORT uint = 8000

var MigrationInstance interfaces.IMigrations

var Configure func()

func UseEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		panic(err.Error())
	}
	ENV_PATH = path
}

// GetEnv returns the value of the environment variable or a default value if not set.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
