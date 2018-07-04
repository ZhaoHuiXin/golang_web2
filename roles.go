package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.autoforce.net/autoforce/admin/model"
	"github.com/gorilla/mux"
)

func init() {
	app.HandleFunc("GET", "/departments/roles", departmentsRolesHandler)
	app.HandleFunc("GET", "/department/{depid:[0-9]+}/roles", rolesHandler)
	app.HandleFunc("POST", "/department/{depid:[0-9]+}/role", roleCreateHandler)
	app.HandleFunc("PUT", "/department/{depid:[0-9]+}/role", nameNotin(updateTableByPKHandler(`roles`), "superior"))
}

func departmentsRolesHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	deps, _ := app.Query("SELECT id,name,superior FROM departments")
	data["departments"] = deps
	roles, _ := app.Query("SELECT id,name,superior FROM roles")
	data["roles"] = roles
	jsonp(w, r, data)
}

func rolesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depid, _ := strconv.ParseInt(vars["depid"], 10, 64)
	data := make(map[string]interface{})
	data["depid"] = vars["depid"]
	var depname string
	err := app.db.QueryRow("SELECT name FROM departments WHERE id=?", depid).Scan(&depname)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data["depname"] = depname

	var roles []model.Role
	rows, err := app.db.Query("SELECT id,name FROM roles WHERE superior=?", depid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name sql.NullString
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			log.Println(err)
			continue
		}
		role := model.Role{
			Id:   id,
			Name: name.String,
		}
		roles = append(roles, role)
	}
	checkRowsError(rows)
	data["roles"] = roles
	data["features"] = app.features
	_locals(r, data, true, true).Render(w, "roles.html")
}

func roleCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depid, _ := strconv.ParseInt(vars["depid"], 10, 64)
	if r.FormValue("superior") == "" {
		r.Form.Set("superior", vars["depid"])
	}

	required := []string{
		"name",
		"superior",
	}

	columns := []string{
		"name",
		"superior",
	}

	table := "`roles`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert roles fail", err)
		redirect_back(w, r, fmt.Sprintf("/department/%d/roles", depid))
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert roles fail", err)
	}
	redirect_back(w, r, fmt.Sprintf("/department/%d/roles", depid))
}
