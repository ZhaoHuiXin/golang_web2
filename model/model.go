package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	table string
	db    string
	cols  []columndef
}

type _Base struct {
	config *Config
	db     *sql.DB
}

var (
	ErrInvalidArg = errors.New("Base: invalid arg")
)

type Model interface {
	Update(pk int64, arg sql.NamedArg) (sql.Result, error)
	Delete(pk int64) (sql.Result, error)
	ToMigrationSql() string
	FullTableName() string
}

type ModelNewFunc func(db *sql.DB) Model

var (
	globalDb *sql.DB
	models   = make(map[string]ModelNewFunc)
)

func RegisterDB(db *sql.DB) {
	if db == nil {
		panic("model: RegisterDB db is nil")
	}
	globalDb = db
}

func Register(name string, newFunc ModelNewFunc) {
	if newFunc == nil {
		panic("model: Register model is nil")
	}
	if _, dup := models[name]; dup {
		panic("model: Register called twice for model " + name)
	}
	models[name] = newFunc
}

func unregisterAllDrivers() {
	// For tests.
	models = make(map[string]ModelNewFunc)
}

func GetModel(key string) Model {
	if newFunc, has := models[key]; has {
		return newFunc(globalDb)
	}
	return nil
}

func (p *_Base) FullTableName() string {
	if p.config.db != "" {
		return "`" + p.config.db + "`.`" + p.config.table + "`"
	}
	return "`" + p.config.table + "`"
}

func (p *_Base) ToMigrationSql() string {
	s, err := to_migration_sql(p.db, p.FullTableName(), p.config.cols)
	if err != nil {
		log.Println("toMigration fail", err)
		return ""
	}
	return s
}

func GenMigrationSql() {
	for k, nfunc := range models {
		mod := nfunc(globalDb)
		s := mod.ToMigrationSql()
		if s == "" {
			continue
		}
		fmt.Printf("-- model: %s\n-- table: %s\n%s\n", k, mod.FullTableName(), s)
	}
}

func (p *_Base) Create(columns, required []string, r *http.Request) (sql.Result, error) {
	var values []interface{}
	var holder []string
	for _, k := range required {
		if r.FormValue(k) == "" {
			return nil, fmt.Errorf("required error")
		}
	}
	for _, k := range columns {
		values = append(values, r.FormValue(k))
		holder = append(holder, `?`)
	}
	sqlstr := fmt.Sprintf(
		"INSERT INTO %s(`%s`) VALUES(%s)",
		p.FullTableName(),
		strings.Join(columns, "`,`"),
		strings.Join(holder, ","))

	return p.db.Exec(sqlstr, values...)
}

func (p *_Base) Read() {
	log.Println("model", p.config.table)
}

func (p *_Base) Update(pk int64, arg sql.NamedArg) (sql.Result, error) {
	if len(arg.Name) < 1 {
		return nil, ErrInvalidArg
	}

	if strings.ContainsAny(arg.Name, "`") {
		return nil, ErrInvalidArg
	}

	if arg.Name == `id` {
		return nil, ErrInvalidArg
	}
	res, err := p.db.Exec(
		fmt.Sprintf("UPDATE %s SET `%s`=? WHERE id=?", p.FullTableName(), arg.Name),
		arg.Value, pk)
	if err != nil {
		log.Println("model exec fail", err)
	}
	return res, err
}

func (p *_Base) Delete(pk int64) (sql.Result, error) {
	res, err := p.db.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE id=?", p.FullTableName()),
		pk)
	if err != nil {
		log.Println("model Delete fail", err)
	}
	return res, err
}

func (p *_Base) List() {
	log.Println("model", p.config.table)
}
