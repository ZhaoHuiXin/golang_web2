package model

import (
	"database/sql"
	"log"
)

type Feature struct {
	_Base
	Id       int
	Name     string
	Path     string
	Category string
	Methods  string
	Auth     bool
}

var features_config = &Config{
	db:    "cheshi",
	table: "features",
	cols:  []columndef{},
}

func NewFeature(db *sql.DB) Model {
	p := &Feature{}
	p._Base.db = db
	p._Base.config = features_config
	return p
}

func init() {
	Register("features", NewFeature)
}

func (p *Feature) Sync() {
	var id int
	var auth bool
	if p.db == nil {
		p.db = globalDb
	}
	err := p.db.QueryRow("SELECT id,`auth` FROM features WHERE `methods`=? AND `path`=?",
		p.Methods, p.Path).Scan(&id, &auth)
	if err == nil && id > 0 {
		p.Id = id
		p.Auth = auth
		return
	}

	sqlstr := "INSERT INTO features(`methods`,`path`,`auth`) VALUES(?,?,?)"
	res, err := p.db.Exec(sqlstr, p.Methods, p.Path, p.Auth)
	if err != nil {
		log.Println("insert features fail", p.Methods, p.Path, err)
		return
	}
	newid, err := res.LastInsertId()
	if err != nil {
		log.Println("get LastInsertId fail", err)
		p.Id = int(newid)
		return
	}
}
