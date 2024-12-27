package controllers

import (
	"simpleApplication/internal/apiserver"
	"simpleApplication/internal/apiserver/core/service"
)

type Controller struct {
	service     *service.Service
	application *apiserver.Application
}

func New(application *apiserver.Application) (*Controller, error) {
	service, err := service.New(application)
	if err != nil {
		return nil, err
	}
	return &Controller{
		service:     service,
		application: application,
	}, nil
}
