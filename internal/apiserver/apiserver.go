package apiserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Config struct {
	ServerBindAddress string
	DbConnectString   string
}
type Application struct {
	Config  *Config
	Logger  *log.Logger
	db      *sql.DB
	methods map[string]*Method
}

type Method struct {
	Add    string
	Get    string
	List   string
	Update string
	Remove string
}
type Field struct {
	Name string
	Type string
	Key  bool
}
type Fields []*Field

func (f *Fields) GetKeyField() *Field {
	for index := range *f {
		if (*f)[index].Key {
			return (*f)[index]
		}
	}
	return nil
}

func (f *Fields) FormatListFieldsForAdd() string {
	listFields := ""
	for index := range *f {
		if (*f)[index].Key {
			continue
		}
		if listFields != "" {
			listFields = listFields + ","
		}
		listFields = fmt.Sprintf("%[1]s %[2]s=($1::json #>> '{body,%[2]s}')", listFields, (*f)[index].Name)
	}
	return listFields
}
func (f *Fields) FormatListFieldsForUpdate() string {

	listFields := ""
	for index := range *f {
		if (*f)[index].Key {
			continue
		}
		if listFields != "" {
			listFields = listFields + ","
		}
		listFields = fmt.Sprintf("%[1]s %[2]s=($1::json #>> '{body,%[2]s}')::%[3]s", listFields, (*f)[index].Name, (*f)[index].Type)
	}
	return listFields
}

type Table struct {
	Schema string
	Name   string
}

func GetFields(db *sql.DB, table *Table) (*Fields, error) {

	fields := &Fields{}
	rows, err := db.Query(`
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
		row := new(Field)
		if err := rows.Scan(&row.Name, &row.Type, &row.Key); err != nil {
			return nil, err
		}
		*fields = append(*fields, row)
	}
	return fields, nil
}
func GetTables(db *sql.DB, schemaName string) ([]*Table, error) {
	tables := []*Table{}
	rows, err := db.Query(`
		SELECT table_schema, table_name 
		FROM information_schema.tables
		WHERE table_schema = $1`, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		row := new(Table)
		if err := rows.Scan(&row.Schema, &row.Name); err != nil {
			return nil, err
		}
		tables = append(tables, row)
	}

	return tables, nil
}
func AddMethod(table *Table, fields string) string {
	return fmt.Sprintf(`
		INSERT INTO %[1]s.%[2]s(%[3]s)
		SELECT %[3]s FROM json_populate_record(NULL::%[1]s.%[2]s, ($1::json -> 'body'))
		RETURNING to_json(%[2]s.*)`,
		table.Schema,
		table.Name,
		fields,
	)
}

func UpdateMethod(table *Table, keyField *Field, fields string) string {
	return fmt.Sprintf(`
		UPDATE %[1]s.%[2]s
		SET %[5]s
		WHERE %[3]s = ($1::json ->>'uri')::%[4]s
		RETURNING to_json(%[2]s.*)`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
		fields,
	)
}
func GetMethod(table *Table, keyField *Field) string {
	return fmt.Sprintf(`
		SELECT row_to_json(%[2]s) FROM %[1]s.%[2]s WHERE %[3]s = ($1::json ->>'uri')::%[4]s`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
	)
}
func ListMethod(table *Table) string {
	return fmt.Sprintf(`
		SELECT json_agg(row_to_json(%[2]s)) FROM %[1]s.%[2]s OFFSET ($1::json #>>'{query,offset}')::integer LIMIT ($1::json #>>'{query,limit}')::integer`,
		table.Schema,
		table.Name,
	)
}
func RemoveMethod(table *Table, keyField *Field) string {
	return fmt.Sprintf(`
		DELETE FROM %[1]s.%[2]s WHERE %[3]s = ($1::json ->>'uri')::%[4]s RETURNING to_json(%[2]s.%[3]s)`,
		table.Schema,
		table.Name,
		keyField.Name,
		keyField.Type,
	)
}
func GetMethods(db *sql.DB, schemaName string) (map[string]*Method, error) {
	methods := map[string]*Method{}
	tables, err := GetTables(db, schemaName)
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		fields, err := GetFields(db, table)
		if err != nil {
			return nil, err
		}
		method := &Method{
			Add:    AddMethod(table, fields.FormatListFieldsForAdd()),
			Get:    GetMethod(table, fields.GetKeyField()),
			List:   ListMethod(table),
			Update: UpdateMethod(table, fields.GetKeyField(), fields.FormatListFieldsForUpdate()),
			Remove: RemoveMethod(table, fields.GetKeyField()),
		}
		methods[table.Name] = method
	}

	return methods, nil
}

func Run() {

	config := &Config{
		ServerBindAddress: ":8000",
		DbConnectString:   "host=localhost database=song_library port=5433 sslmode=disable user=postgres password=1234",
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", config.DbConnectString)
	if err != nil {
		logger.Fatal(err.Error())
	}

	methods, err := GetMethods(db, "public")
	if err != nil {
		logger.Fatal(err.Error())
	}

	application := &Application{
		Config:  config,
		db:      db,
		Logger:  logger,
		methods: methods,
	}

	router := gin.Default()
	application.InitBaseRoutes(router)
	if err := router.Run(config.ServerBindAddress); err != nil {
		logger.Fatal(err.Error())
	}

}

func (a *Application) InitBaseRoutes(router *gin.Engine) {

	for key := range a.methods {
		path := fmt.Sprintf("/%s", key)
		router.POST(path, a.add)

		path = fmt.Sprintf("/%s/:id", key)
		router.PUT(path, a.update)

		path = fmt.Sprintf("/%s/:id", key)
		router.GET(path, a.get)

		path = fmt.Sprintf("/%s", key)
		router.GET(path, a.list)

		path = fmt.Sprintf("/%s/:id", key)
		router.DELETE(path, a.remove)
	}

}

func (a *Application) add(ctx *gin.Context) {
	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(400)
		return
	}
	table := strings.Split(ctx.Request.RequestURI, "/")[1]
	query := a.methods[table].Add
	res, err := a.QueryBase(query, &Params{
		Body: &body,
	})
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (a *Application) update(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(400)
		return
	}
	table := strings.Split(ctx.Request.RequestURI, "/")[1]
	query := a.methods[table].Update

	res, err := a.QueryBase(query, &Params{
		Uri:  &uriParam,
		Body: &body,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}

func (a *Application) get(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	table := strings.Split(ctx.Request.RequestURI, "/")[1]
	query := a.methods[table].Get
	res, err := a.QueryBase(query, &Params{
		Uri: &uriParam,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (a *Application) list(ctx *gin.Context) {
	var queryParam any
	if err := ctx.ShouldBindQuery(&queryParam); err != nil {
		ctx.Status(400)
		return
	}
	table := strings.Split(ctx.Request.RequestURI, "/")[1]
	query := a.methods[table].List
	res, err := a.QueryBase(query, &Params{
		Query: &queryParam,
	})
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.JSON(200, res)
}
func (a *Application) remove(ctx *gin.Context) {
	var uriParam any = ctx.Param("id")
	table := strings.Split(ctx.Request.RequestURI, "/")[1]
	query := a.methods[table].Remove
	_, err := a.QueryBase(query, &Params{
		Uri: &uriParam,
	})
	if errors.Is(err, sql.ErrNoRows) {
		ctx.Status(404)
		return
	}
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.Status(204)

}

type Params struct {
	Uri   *any `json:"uri"`
	Query *any `json:"query"`
	Body  *any `json:"body"`
}

func (a *Application) QueryBase(query string, req *Params) (*any, error) {
	j, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	if err := a.db.QueryRow(query, j).Scan(&j); err != nil {
		return nil, err
	}
	res := new(any)
	err = json.Unmarshal(j, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
