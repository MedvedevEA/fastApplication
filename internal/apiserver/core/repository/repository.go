package repository

import "simpleApplication/internal/apiserver/core/model"

type Repository interface {
	ExecuteQuery(path string, req *model.Params) (*any, error)
	GetListQueries(schemaName string) ([]*model.Query, error)
	SetQueryRoutes(queryRoutes map[string]string)
}
