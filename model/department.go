package model

import "database/sql"

type Department struct {
	_Base
	Id       int
	Name     string
	Superior int
}

var departments_config = &Config{
	db:    "cheshi",
	table: "departments",
	cols:  []columndef{},
}

func NewDepartment(db *sql.DB) Model {
	p := &Department{}
	p._Base.db = db
	p._Base.config = departments_config
	return p
}

func init() {
	Register("departments", NewDepartment)
}
