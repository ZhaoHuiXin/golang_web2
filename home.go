package main

import (
	"net/http"
)

const session_name = `user`

func init() {
	app.HandleFunc("GET", "/", homeHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/users", http.StatusFound)
}
