package service

import (
	"simpleApplication/internal/apiserver/application"
	"simpleApplication/internal/apiserver/core/model"
	"simpleApplication/internal/apiserver/core/repository"
	"simpleApplication/internal/apiserver/infrastructure/sqlrepository"
)

type Service struct {
	application *application.Application
	repository  repository.Repository
}

func New(application *application.Application) (*Service, error) {
	repository, err := sqlrepository.New(application)
	if err != nil {
		return nil, err
	}
	return &Service{
		application: application,
		repository:  repository,
	}, nil
}
func (s *Service) ExecuteQuery(path string, req *model.Params) (*any, error) {
	return s.repository.ExecuteQuery(path, req)
}
func (s *Service) GetListQueries(schemaName string) ([]*model.Query, error) {
	return s.repository.GetListQueries(schemaName)
}
func (s *Service) SetQueryRoutes(queryRoutes map[string]string) {
	s.repository.SetQueryRoutes(queryRoutes)
}
