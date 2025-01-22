package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel        string
	LogFileName     string
	BindAddress     string
	DbConnectString string
	WorkDbSchema    string
}

func New() (*Config, error) {
	if err := godotenv.Load("../../configs/fastApplication.env"); err != nil {
		return nil, err
	}
	return &Config{
		LogLevel:        getEnv("FASTAPPLICATION_SERVER_LOG_LEVEL", "debug"),
		LogFileName:     getEnv("FASTAPPLICATION_SERVER_LOG_FILE_NAME", "fastApplication.log"),
		BindAddress:     getEnv("FASTAPPLICATION_SERVER_BIND_ADDRESS", ":8000"),
		DbConnectString: getEnv("FASTAPPLICATION_DATABASE_CONNECT_STRING", "host=localhost database=fast_application port=5432 sslmode=disable user=postgres password=1234"),
		WorkDbSchema:    getEnv("FASTAPPLICATION_WORK_DATABASE_SCHEMA", "public"),
	}, nil
}
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
