package model

type Table struct {
	Schema string
	Name   string
}
type Field struct {
	Name string
	Type string
	Key  bool
}

type Query struct {
	TableName     string
	BaseQueryName string
	Query         string
}

type Params struct {
	Uri   *any `json:"uri"`
	Query *any `json:"query"`
	Body  *any `json:"body"`
}
