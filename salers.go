package main

import (
	"database/sql"
	"log"
	"net/http"

	"git.autoforce.net/autoforce/admin/model"
)

func init() {
	app.HandleFunc("GET", "/ixiao/salers", salersHandler)
	app.HandleFunc("POST", "/ixiao/saler", salerCreateHandler)
	app.HandleFunc("PUT", "/ixiao/saler", updateTableByPKHandler(`cheyixiao.salers`))
	app.HandleFunc("GET", "/ixiao/saler", salerHandler)
}

func salersHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}
	rows, err := app.db.Query("SELECT id,username,name,phone,cert_code FROM cheyixiao.salers WHERE id>? LIMIT 0,30", start)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var salers []model.Saler
	for rows.Next() {
		var id int
		var name sql.NullString
		var username sql.NullString
		var code sql.NullInt64
		var phone sql.NullInt64
		if err := rows.Scan(&id, &username, &name, &phone, &code); err != nil {
			log.Fatal(err)
		}
		saler := model.Saler{
			Id:       id,
			Username: username.String,
			Name:     name.String,
			Phone:    phone.Int64,
			CertCode: int(code.Int64),
		}
		salers = append(salers, saler)
	}
	checkRowsError(rows)
	data["salers"] = salers
	_locals(r, data, true, true).Render(w, "salers.html")
}

func salerCreateHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"username",
		"name",
	}

	columns := []string{
		"username",
		"name",
		"phone",
	}

	table := "cheyixiao.salers"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert salers fail", err)
		redirect_back(w, r, "/ixiao/salers")
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert salers fail", err)
	}
	redirect_back(w, r, "/ixiao/salers")
}

func salerHandler(w http.ResponseWriter, r *http.Request) {
}
