package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// GetEnv returns the value of the environment variable or a default value if not set.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

var (
	ENVIRONMENT string
	// Common
	SERVICE_NAME     string
	SERVICE_HOST     string
	SERVICE_PORT     string
	SERVICE_TZ       string

	// Postgre
	DATABASE_POSTGRESQL_HOST     string
	DATABASE_POSTGRESQL_PORT     int
	DATABASE_POSTGRESQL_USER     string
	DATABASE_POSTGRESQL_PASSWORD string
	DATABASE_POSTGRESQL_DB_NAME  string

)

func Configure() {
	//try load env from dotenv file if file exists
	godotenv.Load(".env")
	ENVIRONMENT = GetEnv("ENVIRONMENT", "production")

	SERVICE_TZ = GetEnv("SERVICE_TZ", "Asia/Jakarta")

	DATABASE_POSTGRESQL_HOST = GetEnv("DATABASE_POSTGRESQL_HOST", "localhost")
	DATABASE_POSTGRESQL_PORT, _ = strconv.Atoi(GetEnv("DATABASE_POSTGRESQL_PORT", "443"))
	DATABASE_POSTGRESQL_USER = GetEnv("DATABASE_POSTGRESQL_USER", "root")
	DATABASE_POSTGRESQL_PASSWORD = GetEnv("DATABASE_POSTGRESQL_PASSWORD", "root")
	DATABASE_POSTGRESQL_DB_NAME = GetEnv("DATABASE_POSTGRESQL_DB_NAME", "postgres")
}
