package sqlrepository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"simpleApplication/internal/apiserver/application"
	"simpleApplication/internal/apiserver/core/model"

	_ "github.com/lib/pq"
)

type SqlRepository struct {
	application *application.Application
	db          *sql.DB
	QueryRoutes map[string]string
}

func New(application *application.Application) (*SqlRepository, error) {
	db, err := sql.Open("postgres", application.Config.DbConnectString)
	if err != nil {
		return nil, err
	}
	return &SqlRepository{
		application: application,
		db:          db,
		QueryRoutes: nil,
	}, nil
}

func (sr *SqlRepository) ExecuteQuery(path string, req *model.Params) (*any, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	row := sr.db.QueryRow(sr.QueryRoutes[path], j)
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
func (sr *SqlRepository) GetListQueries(schemaName string) ([]*model.Query, error) {
	baseQueries := map[string]func(table *model.Table, fields []*model.Field) string{
		"add":    GetQueryAdd,
		"get":    GetQueryGet,
		"list":   GetQueryList,
		"update": GetQueryUpdate,
		"remove": GetQueryRemove,
	}
	tables, err := sr.GetListTables(schemaName)
	if err != nil {
		return nil, err
	}
	listQueries := []*model.Query{}

	for _, table := range tables {
		fields, err := sr.GetListFields(table)
		if err != nil {
			return nil, err
		}
		for baseQueryKey, baseQuery := range baseQueries {
			listQueries = append(listQueries, &model.Query{
				TableName:     table.Name,
				BaseQueryName: baseQueryKey,
				Query:         baseQuery(table, fields),
			})
		}
	}
	return listQueries, nil
}
func (sr *SqlRepository) SetQueryRoutes(queryRoutes map[string]string) {
	sr.QueryRoutes = queryRoutes
}
func (sr *SqlRepository) GetListTables(schemaName string) ([]*model.Table, error) {
	tables := []*model.Table{}
	rows, err := sr.db.Query(`
		SELECT table_schema, table_name 
		FROM information_schema.tables
		WHERE table_schema = $1`, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		row := new(model.Table)
		if err := rows.Scan(&row.Schema, &row.Name); err != nil {
			return nil, err
		}
		tables = append(tables, row)
	}
	return tables, nil
}
func (sr *SqlRepository) GetListFields(table *model.Table) ([]*model.Field, error) {
	fields := []*model.Field{}
	rows, err := sr.db.Query(`
		SELECT
			_c.column_name ,
			_c.data_type,
			_tc.constraint_type IS NOT NULL
		FROM information_schema.columns _c
		LEFT JOIN information_schema.constraint_column_usage _ccu
		ON
			_c.table_schema=_ccu.table_schema AND
			_c.table_name=_ccu.table_name AND
			_c.column_name=_ccu.column_name
		LEFT JOIN information_schema.table_constraints _tc
		ON
			_tc.constraint_name=_ccu.constraint_name AND
			_tc.constraint_type='PRIMARY KEY'
		WHERE
			_c.table_schema=$1 AND
			_c.table_name=$2`, table.Schema, table.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		row := new(model.Field)
		if err := rows.Scan(&row.Name, &row.Type, &row.Key); err != nil {
			return nil, err
		}
		fields = append(fields, row)
	}
	return fields, nil
}

func GetQueryAdd(table *model.Table, fields []*model.Field) string {
	listFields := ""
	for index := range fields {
		if fields[index].Key {
			continue
		}
		if listFields != "" {
			listFields = listFields + ","
		}
		listFields = fmt.Sprintf("%[1]s %[2]s=($1::json #>> '{body,%[2]s}')", listFields, fields[index].Name)
	}
	return fmt.Sprintf(`
		INSERT INTO %[1]s.%[2]s(%[3]s)
		SELECT %[3]s FROM json_populate_record(NULL::%[1]s.%[2]s, ($1::json -> 'body'))
		RETURNING to_json(%[2]s.*)`,
		table.Schema,
		table.Name,
		listFields,
	)
}

func GetQueryUpdate(table *model.Table, fields []*model.Field) string {
	var listFields string
	var keyField *model.Field
	for index := range fields {
		if fields[index].Key {
			keyField = fields[index]
		}
		if listFields != "" {
			listFields = listFields + ","
		}
		listFields = fmt.Sprintf("%[1]s %[2]s=($1::json #>> '{body,%[2]s}')::%[3]s", listFields, fields[index].Name, fields[index].Type)
	}
	return fmt.Sprintf(`
		UPDATE %[1]s.%[2]s
		SET %[5]s
		WHERE %[3]s = ($1::json ->>'uri')::%[4]s
		RETURNING to_json(%[2]s.*)`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
		listFields,
	)
}

func GetQueryGet(table *model.Table, fields []*model.Field) string {
	var keyField *model.Field
	for index := range fields {
		if fields[index].Key {
			keyField = fields[index]
			break
		}
	}
	return fmt.Sprintf(`
		SELECT row_to_json(%[2]s) FROM %[1]s.%[2]s WHERE %[3]s = ($1::json ->>'uri')::%[4]s`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
	)
}
func GetQueryList(table *model.Table, fields []*model.Field) string {
	return fmt.Sprintf(`
		SELECT json_agg(row_to_json(%[2]s)) FROM %[1]s.%[2]s OFFSET ($1::json #>>'{query,offset}')::integer LIMIT ($1::json #>>'{query,limit}')::integer`,
		table.Schema,
		table.Name,
	)
}

func GetQueryRemove(table *model.Table, fields []*model.Field) string {
	var keyField *model.Field
	for index := range fields {
		if fields[index].Key {
			keyField = fields[index]
			break
		}
	}
	return fmt.Sprintf(`
		DELETE FROM %[1]s.%[2]s WHERE %[3]s = ($1::json ->>'uri')::%[4]s RETURNING to_json(%[2]s.%[3]s)`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
	)
}
