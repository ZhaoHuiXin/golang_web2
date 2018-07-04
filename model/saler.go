package model

import (
	"database/sql"
)

type Saler struct {
	_Base
	Id        int
	Username  string
	Name      string
	Avatar    string
	Phone     int64
	Address   string
	CertCode  int
	UpdatedAt string
	Role      int64
	Brands    string
	IDFace    string
	IDCon     string
	City      string
	DealerId  int64
	CityName  string
}

var salers_config = &Config{
	db:    "cheyixiao",
	table: "salers",
	cols:  []columndef{},
}

func NewSaler(db *sql.DB) Model {
	p := &Saler{}
	p._Base.db = db
	p._Base.config = salers_config
	return p
}

func init() {
	Register("cheyixiao.salers", NewSaler)
}
