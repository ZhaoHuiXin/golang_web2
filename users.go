package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"git.autoforce.net/autoforce/admin/model"
)

func init() {
	app.HandleFunc("GET", "/users", usersHandler)
	app.HandleFunc("POST", "/user", userHandler)
	app.HandleFunc("PUT", "/user", updateTableByPKHandler(`users`))
	app.HandleFunc("GET", "/profile", profileHandler)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	_locals(r, nil, true, true).Render(w, "profile.html")
}

func isWbind(unionid string) bool {
	return strings.HasPrefix(unionid, "waitbind:")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}
	rows, err := app.db.Query("SELECT a.id,a.serial,a.name,a.email,a.unionid,a.dep_id,a.role_id,b.name as dep,c.name as role FROM cheshi.users a LEFT JOIN cheshi.departments b ON b.id=a.dep_id LEFT JOIN cheshi.roles c ON c.id=a.role_id WHERE a.id>? LIMIT 0,30", start)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var serial sql.NullString
		var email sql.NullString
		var unionid sql.NullString
		var id int
		var name sql.NullString
		var dep_id sql.NullInt64
		var role_id sql.NullInt64
		var dep sql.NullString
		var role sql.NullString
		if err := rows.Scan(&id, &serial, &name, &email, &unionid, &dep_id, &role_id, &dep, &role); err != nil {
			log.Fatal(err)
		}
		user := model.User{
			Id:     id,
			Serial: serial.String,
			Name:   name.String,
			Email:  email.String,
			Wbind:  isWbind(unionid.String),
			DepId:  int(dep_id.Int64),
			RoleId: int(role_id.Int64),
			Dep:    dep.String,
			Role:   role.String,
		}
		if user.Wbind {
			user.QRCode = fmt.Sprintf(
				"https://api.autoforce.net/wx/dms/qrcode?action_name=QR_STR_SCENE&scene=2:%d,0&expire_seconds=120",
				id)
		}
		users = append(users, user)
	}
	checkRowsError(rows)
	data["users"] = users
	_locals(r, data, true, true).Render(w, "users.html")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		return
	}

	serial := r.FormValue("serial")
	name := r.FormValue("name")
	email := r.FormValue("email")
	dep_id := r.FormValue("dep_id")
	role_id := r.FormValue("role_id")
	num := rand.Intn(100)
	unionid := `waitbind:` + strconv.Itoa(num)
	res, err := app.db.Exec(
		`INSERT INTO cheshi.users(serial,name,email,unionid,dep_id,role_id) VALUES(?,?,?,?,?,?)`,
		serial, name, email, unionid, dep_id, role_id)
	if err != nil {
		session.AddFlash(`err:` + err.Error())
		session.Save(r, w)
		log.Println("insert user fail", serial, name, email, err)
		redirect_back(w, r, "/users")
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		session.AddFlash(`err:` + err.Error())
		session.Save(r, w)
		log.Println("insert user fail", serial, name, email, err)
		redirect_back(w, r, "/users")
		return
	}
	session.AddFlash(fmt.Sprintf("newuser:%d", id))
	session.Save(r, w)
	redirect_back(w, r, "/users")
}
