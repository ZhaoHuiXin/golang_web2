package main

import (
	"log"
	"net/http"
)

func init() {
	app.HandleFunc("GET", "/features", featuresHandler)
	app.HandleFunc("PUT", "/feature", updateTableByPKHandler(`features`))
	app.HandleFunc("DELETE", "/feature", deleteTableByPKHandler(`features`))
}

func featuresHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	res, err := app.Query("SELECT id,name,path,methods,auth FROM features")
	if err != nil {
		log.Fatal(err)
	}
	data["features"] = res
	_locals(r, data, true, true).Render(w, "features.html")
}
