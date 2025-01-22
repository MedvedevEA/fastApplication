package main

import (
	"fastApplication/internal/apiserver"
	"fastApplication/internal/config"
	"fastApplication/internal/logger/logrus"
	"fastApplication/internal/service"
	"fastApplication/internal/sqlrepository"
	"log"
	"os"
)

func main() {
	//config
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	//logger
	logFile, err := os.OpenFile(config.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logger, err := logrus.New(config.LogLevel, os.Stdout, logFile)
	if err != nil {
		log.Fatal(err)
	}
	//store
	store, err := sqlrepository.New(config.DbConnectString, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	//service
	service := service.New(store, logger)
	//apiserver
	logger.Infof("API Server 'Fast Application' is started in addr:[%s]", config.BindAddress)
	apiServer := apiserver.New(config.BindAddress, service, logger, config.WorkDbSchema)
	if err := apiServer.Run(); err != nil {
		logger.Errorf("API Server 'Fast Application' error: %s", err)
		return
	}
	logger.Infof("API Server 'Fast Application' is stoped")

}
