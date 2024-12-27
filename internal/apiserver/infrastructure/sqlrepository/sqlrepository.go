package sqlrepository

import (
	"database/sql"
	"encoding/json"
	"simpleApplication/internal/apiserver"

	_ "github.com/lib/pq"
)

type SqlRepository struct {
	application *apiserver.Application
	db          *sql.DB
}

func New(application *apiserver.Application) (*SqlRepository, error) {
	db, err := sql.Open("postgres", application.Config.DbConnectString)
	if err != nil {
		return nil, err
	}
	return &SqlRepository{
		application: application,
		db:          db,
	}, nil
}

func (s *SqlRepository) ExecuteSqlQuery(query string, req *any) (*any, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	row := s.db.QueryRow(query, j)
	if err := row.Scan(&j); err != nil {
		return nil, err
	}
	res := new(any)
	err = json.Unmarshal(j, res)
	if err != nil {
		return nil, err
	}
	return res, nil

}
