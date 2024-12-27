package apiserver

import (
	"simpleApplication/internal/apiserver/application"
	"simpleApplication/internal/apiserver/infrastructure/controller"

	"github.com/gin-gonic/gin"
)

func Run() {
	application := application.New()

	router := gin.Default()

	if err := controller.Init(application, router); err != nil {
		application.Logger.Fatal(err.Error())
	}
	if err := router.Run(application.Config.ServerBindAddress); err != nil {
		application.Logger.Fatal(err.Error())
	}

}
