package app

import (
	"fastApplication/internal/controller"
	"fastApplication/internal/service"
	"fastApplication/internal/sqlrepository"
	"log"

	"github.com/gin-gonic/gin"
)

type Config struct {
	ServerBindAddress string
	DbConnectString   string
}

func Run() {
	config := &Config{
		ServerBindAddress: ":8000",
		DbConnectString:   "host=localhost database=song_library port=5433 sslmode=disable user=postgres password=1234",
	}

	r, err := sqlrepository.New(config.DbConnectString)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	s := service.New(r)

	router := gin.Default()
	controller.RegisterRoutes(router, s, "public")

	if err := router.Run(config.ServerBindAddress); err != nil {
		log.Fatal(err.Error())
	}

}
