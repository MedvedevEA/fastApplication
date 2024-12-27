package service

import (
	"simpleApplication/internal/apiserver"
	"simpleApplication/internal/apiserver/core/repository"
	"simpleApplication/internal/apiserver/infrastructure/sqlrepository"
)

type Service struct {
	application *apiserver.Application
	repository  repository.Repository
}

func New(application *apiserver.Application) (*Service, error) {
	repository, err := sqlrepository.New(application)
	if err != nil {
		return nil, err
	}
	return &Service{
		application: application,
		repository:  repository,
	}, nil
}

/*
func New() *Service {
	sqlRepository := new(sqlrepository.SqlRepository)

	return &Service{
		repository: sqlRepository,
	}

}

func (s *Service) ExecuteSqlQuery(query string, req *any) (*any, error) {
	return s.repository.ExecuteSqlQuery(query, req)
}
*/
