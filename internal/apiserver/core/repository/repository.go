package repository

type Repository interface {
	ExecuteSqlQuery(query string, req *any) (*any, error)
}
