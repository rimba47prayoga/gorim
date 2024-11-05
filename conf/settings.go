package conf

import (
	"os"

	"gorm.io/gorm"
)

var DB *gorm.DB

type Database struct {
	Name		string
	Host		string
	Port		int
	User		string
	Password	string
}


// GetEnv returns the value of the environment variable or a default value if not set.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
