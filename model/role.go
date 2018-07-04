package model

import "database/sql"

type Role struct {
	_Base
	Id       int
	Name     string
	Superior int
}

var roles_config = &Config{
	db:    "cheshi",
	table: "roles",
	cols:  []columndef{},
}

func NewRole(db *sql.DB) Model {
	p := &Role{}
	p._Base.db = db
	p._Base.config = roles_config
	return p
}

func init() {
	Register("roles", NewRole)
}
