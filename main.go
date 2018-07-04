package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"git.autoforce.net/autoforce/admin/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

var render Render

const hashKey = `qwoerxcvijqn3m23*(*Yasdf)asf!@#afasdflkq3r09uijxcv3423424444442`
const sessionDir = `sessions`

var session_store = sessions.NewFilesystemStore(sessionDir, []byte(hashKey))

type Opt struct {
	debug   bool
	https   bool
	migrate bool
	listen  string
}

var opt Opt

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&opt.debug, "debug", false, "debug")
	flag.BoolVar(&opt.https, "https", false, "https")
	flag.BoolVar(&opt.migrate, "migrate", false, "db migrate")
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if port < 1 {
		port = 3000
	}
	flag.StringVar(&opt.listen, "listen", fmt.Sprint(":", port), "listen in")
	os.MkdirAll(sessionDir, 0755)
}

func serve() {
	app.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("serve on", opt.listen)
	loggedRouter := handlers.LoggingHandler(os.Stdout, handlers.RecoveryHandler()(app.router))
	log.Fatalln(http.ListenAndServe(opt.listen, loggedRouter))
}

func main() {
	flag.Parse()
	if !opt.debug {
		opt.debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	}
	log.Printf("debug %v", opt.debug)

	app.Init()

	if opt.migrate {
		model.GenMigrationSql()
		fmt.Printf("migrate -path migrations/cheshi -database '%s' -verbose up\n", app.db_dsn)
		return
	}
	render = NewRender("template", opt.debug)

	serve()
}
