package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	app.HandleFunc("GET", "/role/{role_id:[0-9]+}/access", roleAccessHandler)
	app.HandleFunc("PUT", "/role/{role_id:[0-9]+}/access", updateRoleAccessHandler)
	app.HandleFunc("GET", "/user/{user_id:[0-9]+}/access", userAccessHandler)
	app.HandleFunc("PUT", "/user/{user_id:[0-9]+}/access", updateUserAccessHandler)
}

func roleAccessHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	vars := mux.Vars(r)
	role_id, _ := strconv.ParseInt(vars["role_id"], 10, 64)

	rows, err := app.db.Query("SELECT feat_id FROM cheshi.access where role_id=?", role_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var feats []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		feats = append(feats, id)
	}
	checkRowsError(rows)
	data["feats"] = feats
	data["role_id"] = role_id
	jsonp(w, r, data)
}

func updateRoleAccessHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	role_id, _ := strconv.ParseInt(vars["role_id"], 10, 64)
	if role_id == 1 {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err := app.db.Exec("DELETE FROM cheshi.access where role_id=?", role_id)
	if err != nil {
		log.Println("removeAll access for role ", role_id, err)
		return
	}

	r.FormValue("value[]")
	for _, vs := range r.Form["value[]"] {
		sqlstr := "INSERT INTO cheshi.access(`role_id`,`feat_id`) VALUES(?,?)"
		_, err := app.db.Exec(sqlstr, role_id, vs)
		if err != nil {
			log.Println("insert access fail", err)
		}
	}
	app.LoadAccess()
}

func userAccessHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	vars := mux.Vars(r)
	user_id, _ := strconv.ParseInt(vars["user_id"], 10, 64)

	rows, err := app.db.Query("SELECT feat_id FROM cheshi.access where user_id=?", user_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var feats []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		feats = append(feats, id)
	}
	checkRowsError(rows)
	data["feats"] = feats
	data["user_id"] = user_id
	jsonp(w, r, data)
}

func updateUserAccessHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id, _ := strconv.ParseInt(vars["user_id"], 10, 64)

	_, err := app.db.Exec("DELETE FROM cheshi.access where user_id=?", user_id)
	if err != nil {
		log.Println("removeAll access for user ", user_id, err)
		return
	}

	r.FormValue("value[]")
	for _, vs := range r.Form["value[]"] {
		sqlstr := "INSERT INTO cheshi.access(`user_id`,`feat_id`) VALUES(?,?)"
		_, err := app.db.Exec(sqlstr, user_id, vs)
		if err != nil {
			log.Println("insert access fail", err)
		}
	}
	app.LoadAccess()
}
