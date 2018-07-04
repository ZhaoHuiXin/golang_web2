package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"git.autoforce.net/autoforce/admin/model"
)

func init() {
	app.HandleFunc("GET", "/projects", projectsHandler)
	app.HandleFunc("POST", "/project", projectCreateHandler)
	app.HandleFunc("PUT", "/project", updateTableByPKHandler(`projects`))
	app.HandleFunc("GET", "/project/{id:[0-9]+}", projectHandler)
	app.HandleFunc("POST", "/project/{id:[0-9]+}/staff", projectAddStaffHandler)
}

func projectAddStaffHandler(w http.ResponseWriter, r *http.Request) {
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}
	rows, err := app.db.Query("SELECT id,name,parent_id,leader_id,created_at,archived_at FROM projects WHERE id>? LIMIT 0,30", start)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var projects []model.Project
	for rows.Next() {
		var name sql.NullString
		var parentId sql.NullInt64
		var leaderId sql.NullInt64
		var createdAt sql.NullInt64
		var archivedAt sql.NullInt64
		var id int
		if err := rows.Scan(&id, &name, &parentId, &leaderId, &createdAt, &archivedAt); err != nil {
			log.Fatal(err)
		}
		project := model.Project{
			Id:         id,
			Name:       name.String,
			ParentId:   int(parentId.Int64),
			LeaderId:   int(leaderId.Int64),
			CreatedAt:  int(createdAt.Int64),
			ArchivedAt: int(archivedAt.Int64),
		}
		projects = append(projects, project)
	}
	checkRowsError(rows)
	data["projects"] = projects
	_locals(r, data, true, true).Render(w, "projects.html")
}

func projectCreateHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"name",
	}

	columns := []string{
		"name",
	}

	table := "`projects`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert projects fail", err)
		redirect_back(w, r, "/projects")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("insert projects fail", err)
		redirect_back(w, r, "/projects")
		return
	}
	redirect_back(w, r, fmt.Sprint("/project/", id))
}
