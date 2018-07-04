package model

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type columndef struct {
	Name, Type    string
	PK, notNull   bool
	Key, default_ NullString
	extra         string

	auto_increment bool
}

var (
	IdDef        = UintField("id", 11, "0").AsPK().AsAutoIncrement()
	UpdatedAtDef = TimestampField("updated_at", "CURRENT_TIMESTAMP").Extra("on update CURRENT_TIMESTAMP")
	CreatedAtDef = TimestampField("created_at", "CURRENT_TIMESTAMP")
	DeletedAtDef = NullTimestampField("deleted_at")
)

func BigIntField(name string, length int, defv string) columndef {
	return columndef{
		Name:     name,
		Type:     fmt.Sprintf("BIGINT(%d)", length),
		notNull:  true,
		default_: ToNullString(defv),
	}
}

func NullBigIntField(name string, length int) columndef {
	return columndef{
		Name: name,
		Type: fmt.Sprintf("BIGINT(%d)", length),
	}
}

func IntField(name string, length int, defv string) columndef {
	return columndef{
		Name:     name,
		Type:     fmt.Sprintf("INT(%d)", length),
		notNull:  true,
		default_: ToNullString(defv),
	}
}

func NullIntField(name string, length int) columndef {
	return columndef{
		Name: name,
		Type: fmt.Sprintf("INT(%d)", length),
	}
}

func UintField(name string, length int, defv string) columndef {
	return columndef{
		Name:     name,
		Type:     fmt.Sprintf("INT(%d) UNSIGNED", length),
		notNull:  true,
		default_: ToNullString(defv),
	}
}

func NullUintField(name string, length int) columndef {
	return columndef{
		Name: name,
		Type: fmt.Sprintf("INT(%d) UNSIGNED", length),
	}
}

func TimestampField(name, defv string) columndef {
	return columndef{
		Name:     name,
		Type:     "TIMESTAMP",
		notNull:  true,
		default_: ToNullString(defv),
	}
}

func NullTimestampField(name string) columndef {
	return columndef{
		Name: name,
		Type: "TIMESTAMP",
	}
}

func VarcharField(name string, length int, defv string) columndef {
	return columndef{
		Name:     name,
		Type:     fmt.Sprintf("VARCHAR(%d)", length),
		notNull:  true,
		default_: ToNullString(defv),
	}
}

func NullVarcharField(name string, length int) columndef {
	return columndef{
		Name: name,
		Type: fmt.Sprintf("VARCHAR(%d)", length),
	}
}

func TextField(name string) columndef {
	return columndef{
		Name: name,
		Type: `TEXT`,
	}
}

func (d columndef) AsPK() columndef {
	d.PK = true
	return d
}

func (d columndef) AsAutoIncrement() columndef {
	d.auto_increment = true
	return d
}

func (d columndef) Default(v string) columndef {
	d.default_ = ToNullString(v)
	return d
}

func (d columndef) Extra(extra string) columndef {
	d.extra = extra
	return d
}

func (p *columndef) to_column_sql() string {
	s := "`" + p.Name + "` " + p.Type
	if strings.ToUpper(p.Type) == "TEXT" {
		return s
	}
	if p.notNull || p.PK {
		s += " NOT NULL"
	}
	if p.auto_increment {
		s += " AUTO_INCREMENT"
	} else {
		if p.default_.Valid {
			s += " DEFAULT " + p.default_.String
		}
	}
	if len(p.extra) > 0 {
		s += " " + p.extra
	}
	return s
}

func (p *columndef) sameas(t *columndef) bool {
	same := p.Name == t.Name && strings.ToLower(p.Type) == strings.ToLower(t.Type) && p.notNull == t.notNull && p.PK == t.PK && p.extra == t.extra && p.auto_increment == t.auto_increment
	if !same {
		return same
	}
	if !p.auto_increment {
		return p.default_ == t.default_
	}
	return same
}

func (p *columndef) to_add_column_sql(table string) string {
	return `ALTER TABLE ` + table + ` ADD COLUMN ` + p.to_column_sql() + ";\n"
}

func (p *columndef) to_modify_column_sql(table string) string {
	return `ALTER TABLE ` + table + ` MODIFY COLUMN ` + p.to_column_sql() + ";\n"
}

func to_create_table_sql(table string, columns []columndef) string {
	var pks []string
	for _, cdef := range columns {
		if cdef.PK {
			pks = append(pks, cdef.Name)
		}
	}

	s := `CREATE TABLE IF NOT EXISTS ` + table + "(\n"
	for _, cdef := range columns {
		s += cdef.to_column_sql()
		if cdef.PK && len(pks) == 1 {
			s += " PRIMARY KEY"
		}
		s += ",\n"
	}
	if len(pks) > 1 {
		s += "PRIMARY KEY(" + strings.Join(pks, ", ") + "),\n"
	}
	s = strings.TrimSuffix(s, ",\n")
	s += "\n);\n"
	return s
}

func exists_table(db *sql.DB, table string) bool {
	count := -1
	db.QueryRow("SELECT COUNT(1) FROM " + table + " LIMIT 1").Scan(&count)
	return count > -1
}

func to_migration_sql(db *sql.DB, table string, columns []columndef) (string, error) {
	if !exists_table(db, table) {
		return to_create_table_sql(table, columns), nil
	}

	rows, err := db.Query("DESC " + table)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer rows.Close()

	var onbeach []columndef
	for rows.Next() {
		var t columndef
		var NULL string
		if err := rows.Scan(&t.Name, &t.Type, &NULL, &t.Key, &t.default_, &t.extra); err != nil {
			log.Fatalln(err)
			continue
		}
		extra := strings.Fields(t.extra)
		remove_i := -1
		auto_increment := false
		for i, s := range extra {
			if strings.ToUpper(s) == "AUTO_INCREMENT" {
				auto_increment = true
				remove_i = i
				break
			}
		}
		if remove_i > -1 {
			i := remove_i
			extra = append(extra[:i], extra[i+1:]...)
			t.extra = strings.Join(extra, " ")
		}
		t.auto_increment = auto_increment
		t.notNull = strings.HasPrefix(strings.ToUpper(NULL), "N")
		t.PK = strings.ToUpper(t.Key.String) == "PRI"

		onbeach = append(onbeach, t)
	}
	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	var b bytes.Buffer
	for _, t := range columns {
		ismiss := true
		ischange := false
		for _, cdef := range onbeach {
			if cdef.Name == t.Name {
				ismiss = false
				ischange = !cdef.sameas(&t)
				break
			}
		}
		if ismiss {
			s := t.to_add_column_sql(table)
			b.WriteString(s)
		} else if ischange {
			s := t.to_modify_column_sql(table)
			b.WriteString(s)
		}
	}

	return b.String(), nil
}
