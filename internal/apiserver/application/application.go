package application

import (
	"log"
	"os"
)

type Config struct {
	ServerBindAddress string
	DbConnectString   string
}
type Application struct {
	Config *Config
	Logger *log.Logger
}

func New() *Application {
	config := &Config{
		ServerBindAddress: ":8000",
		DbConnectString:   "host=localhost database=song_library port=5433 sslmode=disable user=postgres password=1234",
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return &Application{
		Config: config,
		Logger: logger,
	}
}
