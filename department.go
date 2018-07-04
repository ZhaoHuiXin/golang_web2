package main

import (
	"log"
	"net/http"
)

func init() {
	app.HandleFunc("GET", "/departments", departmentsHandler)
	app.HandleFunc("POST", "/department", departmentCreateHandler)
	app.HandleFunc("PUT", "/department", updateTableByPKHandler(`departments`))
}

func departmentsHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, err := app.Query("SELECT a.id,a.name,a.superior,b.name as supname FROM departments a LEFT JOIN departments b ON b.id=a.superior")
	if err != nil {
		log.Fatal(err)
	}
	data["departments"] = res
	_locals(r, data, true, true).Render(w, "departments.html")
}

func departmentCreateHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"name",
	}

	columns := []string{
		"name",
		"superior",
	}

	table := "`departments`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert departments fail", err)
		redirect_back(w, r, "/departments")
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert departments fail", err)
	}
	redirect_back(w, r, "/departments")
}
