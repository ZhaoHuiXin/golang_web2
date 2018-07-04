package main

import (
	"log"
	"net/http"
)

func init() {
	app.HandleFunc("GET", "/assets", assetsHandler)
	app.HandleFunc("POST", "/asset", assetHandler)
	app.HandleFunc("PUT", "/asset", updateTableByPKHandler(`assets`))
	app.HandleFunc("GET", "/assets/e", eassetsHandler)
	app.HandleFunc("POST", "/assets/e", eassetHandler)
	app.HandleFunc("PUT", "/assets/e", updateTableByPKHandler(`eassets`))
	app.HandleFunc("GET", "/assets/category", assetsCategoryHandler)
	app.HandleFunc("POST", "/assets/category", assetsCategoryCreateHandler)
	app.HandleFunc("PUT", "/assets/category", updateTableByPKHandler(`assets_category`))
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}
	res, err := app.Query("SELECT id,code,label,kind,model,serial,buy_at,buy_val,comment FROM assets WHERE id>? LIMIT 0,30", start)
	if err != nil {
		log.Fatal(err)
	}
	data["assets"] = res

	res, err = app.Query("SELECT id, name FROM assets_category WHERE parent=1")
	if err != nil {
		log.Fatal(err)
	}
	data["categories"] = res
	_locals(r, data, true, true).Render(w, "assets.html")
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"code",
		"kind",
		"model",
		"serial",
		"buy_at",
		"buy_val",
	}

	columns := []string{
		"label",
		"code",
		"kind",
		"model",
		"serial",
		"buy_at",
		"buy_val",
		"comment",
	}

	table := "`assets`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert asset fail", err)
		redirect_back(w, r, "/assets")
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert asset fail", err)
	}
	redirect_back(w, r, "/assets")
}

func eassetsHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	if start == "" {
		start = "0"
	}
	res, err := app.Query("SELECT id,label,kind,model,serial,buy_at,cost,op,comment FROM eassets WHERE id>? LIMIT 0,30", start)
	if err != nil {
		log.Fatal(err)
	}
	data["eassets"] = res

	res, err = app.Query("SELECT id, name FROM assets_category WHERE parent=2")
	if err != nil {
		log.Fatal(err)
	}
	data["categories"] = res
	_locals(r, data, true, true).Render(w, "eassets.html")
}

func eassetHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"kind",
		"model",
		"serial",
		"buy_at",
		"cost",
	}

	columns := []string{
		"label",
		"kind",
		"model",
		"serial",
		"buy_at",
		"cost",
		"op",
		"comment",
	}

	table := "`eassets`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert eassets fail", err)
		redirect_back(w, r, "/assets/e")
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert eassets fail", err)
	}
	redirect_back(w, r, "/assets/e")
}

func assetsCategoryHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, err := app.Query("SELECT a.id as id, a.name as name, a.parent as parent, b.name as pname FROM assets_category a LEFT JOIN assets_category b ON b.id=a.parent")
	if err != nil {
		log.Fatal(err)
	}
	data["categories"] = res
	_locals(r, data, true, true).Render(w, "assets_category.html")
}

func assetsCategoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{
		"name",
		"parent",
	}

	columns := []string{
		"name",
		"parent",
	}

	table := "`assets_category`"
	res, err := simpleCreateHandle(table, columns, required, r)
	if err != nil {
		log.Println("insert assets_category fail", err)
		redirect_back(w, r, "/assets/category")
		return
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Println("insert assets_category fail", err)
	}
	redirect_back(w, r, "/assets/category")
}
