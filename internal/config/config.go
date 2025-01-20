package config

type Config struct {
	LogLevel                 string
	LogFileName              string
	BindAddress              string
	DatabaseConnectString    string
	OutsideServerBindAddress string
	WorkDbSchema             string
}

func New() (*Config, error) {
	return &Config{
		LogLevel:              "debug",
		LogFileName:           "log.txt",
		BindAddress:           ":8000",
		DatabaseConnectString: "host=localhost database=song_library port=5433 sslmode=disable user=postgres password=1234",
		WorkDbSchema:          "public",
	}, nil
}
