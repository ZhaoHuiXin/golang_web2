package model

import "database/sql"

type Project struct {
	_Base
	Id         int
	Name       string
	ParentId   int
	LeaderId   int
	CreatedAt  int
	ArchivedAt int
}

var projects_config = &Config{
	db:    "cheshi",
	table: "projects",
	cols: []columndef{
		IdDef,
		VarcharField("name", 64, `""`),
		NullIntField("parent_id", 11),
		NullIntField("leader_id", 11),
		NullTimestampField("archived_at"),
		CreatedAtDef,
	},
}

func NewProject(db *sql.DB) Model {
	p := &ProjectMember{}
	p._Base.db = db
	p._Base.config = projects_config
	return p
}

func init() {
	Register("projects", NewProject)
}

type ProjectMember struct {
	_Base
	Id        int
	ProjectId int
	UserId    int
	CreatedAt int
	DeletedAt int
}

var project_member_config = &Config{
	db:    "cheshi",
	table: "project_members",
	cols:  []columndef{},
}

func NewProjectMember(db *sql.DB) Model {
	p := &ProjectMember{}
	p._Base.db = db
	p._Base.config = project_member_config
	return p
}

func init() {
	Register("project_members", NewProjectMember)
}
