package model

import "database/sql"

type User struct {
	_Base
	Id     int
	Serial string
	Name   string
	Email  string
	Wbind  bool
	QRCode string
	DepId  int
	RoleId int
	Dep    string
	Role   string
}

var users_config = &Config{
	db:    "cheshi",
	table: "users",
	cols: []columndef{
		IdDef,
		CreatedAtDef,
		UpdatedAtDef,
		DeletedAtDef,
		NullVarcharField("email", 32),
		NullVarcharField("serial", 32),
	},
}

func NewUser(db *sql.DB) Model {
	p := &User{}
	p._Base.db = db
	p._Base.config = users_config
	return p
}

func init() {
	Register("users", NewUser)
}
