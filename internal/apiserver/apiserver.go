package apiserver

import (
	"context"
	"fastApplication/internal/controller"
	"fastApplication/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	server *http.Server
}

func New(bindAddress string, service controller.Service, logger logger.Logger, workDbSchema string) *ApiServer {
	router := gin.New()
	router.Use(gin.Recovery())
	//router.Use(middlewares.LoggingMiddleware(logger))
	controller.Init(router, service, logger, workDbSchema)

	server := &http.Server{
		Addr:    bindAddress,
		Handler: router,
	}
	return &ApiServer{
		server: server,
	}
}
func (a *ApiServer) Run() error {
	chError := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			chError <- err
		}
	}()
	go func() {
		chQuit := make(chan os.Signal, 1)
		signal.Notify(chQuit, syscall.SIGINT, syscall.SIGTERM)
		<-chQuit
		ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)
		defer channel()
		chError <- a.server.Shutdown(ctx)
	}()

	return <-chError
}
