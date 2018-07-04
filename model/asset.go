package model

import "database/sql"

type Asset struct {
	_Base
	Id      int
	Code    string
	Label   string
	Kind    int
	Model   string
	Serial  string
	BuyAt   string
	BuyVal  string
	Comment string
}

var asset_config = &Config{
	db:    "cheshi",
	table: "assets",
	cols:  []columndef{},
}

func NewAsset(db *sql.DB) Model {
	p := &Asset{}
	p._Base.db = db
	p._Base.config = asset_config
	return p
}

func init() {
	Register("assets", NewAsset)
}
