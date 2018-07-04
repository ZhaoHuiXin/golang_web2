package model

import (
	"database/sql"
	"fmt"
	"log"
)

type Access struct {
	_Base
	RoleId int
	UserId int
	FeatId int
}

var access_config = &Config{
	db:    "cheshi",
	table: "access",
	cols:  []columndef{},
}

func NewAccess(db *sql.DB) Model {
	p := &Access{}
	p._Base.db = db
	p._Base.config = access_config
	return p
}

func init() {
	Register("access", NewAccess)
}

func (p *Access) Remove() (sql.Result, error) {
	res, err := p.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE role_id=? AND user_id=? AND feat_id=?", p.FullTableName()),
		p.RoleId, p.UserId, p.FeatId)
	if err != nil {
		log.Println("access Delete fail", err)
	}
	return res, err
}
