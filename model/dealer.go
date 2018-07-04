package model

import (
	"database/sql"
)

type Dealer struct {
	_Base
	Id         int
	Name       string
	Company    string
	Address    string
	Call       string
	Phone      int64
	ChName     string
	BusLicence string
	PicDoor    string
	PicShow    string
	PicRest    string
	PicOther   string
	Region     string
}

var dealer_config = &Config{
	db:    "cheyixiao",
	table: "dealers",
	cols:  []columndef{},
}

func NewDealer(db *sql.DB) Model {
	p := &Dealer{}
	p._Base.db = db
	p._Base.config = dealer_config
	return p
}

func init() {
	Register("cheyixiao.dealers", NewDealer)
}
