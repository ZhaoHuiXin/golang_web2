package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.autoforce.net/autoforce/admin/model"

	"math"
	"os"

	"github.com/gorilla/sessions"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func redirect_back(w http.ResponseWriter, r *http.Request, defurl string) {
	referer := r.Referer()
	if referer == "" {
		referer = defurl
	}
	http.Redirect(w, r, referer, http.StatusFound)
}

func jsonp(w http.ResponseWriter, r *http.Request, data interface{}) {
	cb := r.FormValue("cb")
	if cb != "" {
		cb = strings.Fields(cb)[0]
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if cb != "" {
		w.Header().Set("Content-Type", "application/javascript; charset=UTF-8")
		fmt.Fprintf(w, "/**/ typeof %s === 'function' && %s(", cb, cb)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	buf, err := json.Marshal(data)
	if err != nil {
		log.Println("marshal data error", err)
	}
	w.Write(buf)
	if cb != "" {
		fmt.Fprintf(w, ");")
	}
}

func expect_json_res(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}
	if strings.Contains(accept, "text/json") {
		return true
	}
	return false
}

func md5str(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

func getSession(r *http.Request) (*sessions.Session, error) {
	return session_store.Get(r, session_name)
}

type MenuItem struct {
	Name    string
	Url     string
	Active  bool
	Submenu *Menu
}

type Menu struct {
	Name   string
	Url    string
	Active bool
	Items  []MenuItem
}

func load_menus(r *http.Request, uid int) (menus []Menu) {
	feats := []Feature{}
	cont, err := ioutil.ReadFile("feat.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(cont, &feats)
	if err != nil {
		log.Println("load feat fail", err)
		return
	}
	for _, f := range feats {
		if f.Category == "" {
			continue
		}
		if f.Category == "self" {
			continue
		}
		fi := -1
		for i, menu := range menus {
			if f.Category == menu.Name {
				fi = i
				break
			}
		}
		active := r.URL.Path == f.Path
		if fi == -1 {
			fi = len(menus)
			menus = append(menus, Menu{
				Name:   f.Category,
				Url:    f.Path,
				Active: active,
			})
		}
		menus[fi].Items = append(menus[fi].Items, MenuItem{
			Name:   f.Name,
			Url:    f.Path,
			Active: active,
		})
		if active {
			menus[fi].Active = true
		}
	}
	return
}

type Locals struct {
	r       *http.Request
	Title   string
	Data    interface{}
	Flashes []interface{}
	Session map[interface{}]interface{}
	Menus   []Menu

	isXhr bool
}

func (p *Locals) IsCurrentPath(u string) bool {
	return p.r.URL.Path == u
}

func (p *Locals) Render(w http.ResponseWriter, page string) {
	if p.isXhr {
		jsonp(w, p.r, p.Data)
		return
	}
	render.Render(w, page, p)
}

func _locals(r *http.Request, data interface{}, loadFlash, loadMenu bool) *Locals {
	locals := &Locals{
		r:     r,
		Data:  data,
		isXhr: r.Header.Get("X-Requested-With") == "XMLHttpRequest",
	}
	session, _ := session_store.Get(r, session_name)
	locals.Session = session.Values

	if loadFlash {
		locals.Flashes = session.Flashes()
	}
	if loadMenu {
		uid := session.Values["uid"]
		if uid != nil {
			locals.Menus = load_menus(r, uid.(int))
		}
	}
	return locals
}

func updateTableByPKHandler(table string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		pk, err := strconv.ParseInt(r.FormValue("pk"), 10, 64)
		if err != nil {
			data["status"] = "fail"
			data["msg"] = err.Error()
			jsonp(w, r, data)
			return
		}
		name := r.FormValue("name")
		value := r.FormValue("value")

		mod := model.GetModel(table)
		if mod == nil {
			log.Println("updateTableByPKHandler: not found model " + table)
			data["status"] = "fail"
			data["msg"] = "no model"
			jsonp(w, r, data)
			return
		}
		_, err = mod.Update(pk, sql.Named(name, value))
		if err != nil {
			data["status"] = "fail"
			data["msg"] = err.Error()
			jsonp(w, r, data)
			return
		}
		data["status"] = "ok"
		jsonp(w, r, data)
		return
	})
}

func simpleCreateHandle(table string, columns, required []string, r *http.Request) (sql.Result, error) {
	var values []interface{}
	var holder []string
	for _, k := range required {
		if r.FormValue(k) == "" {
			return nil, fmt.Errorf("required error")
		}
	}
	for _, k := range columns {
		values = append(values, r.FormValue(k))
		holder = append(holder, `?`)
	}
	sqlstr := fmt.Sprintf(
		"INSERT INTO %s(`%s`) VALUES(%s)",
		table,
		strings.Join(columns, "`,`"),
		strings.Join(holder, ","))

	return app.db.Exec(sqlstr, values...)
}

func deleteTableByPKHandler(table string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		pk, err := strconv.ParseInt(r.FormValue("pk"), 10, 64)
		if err != nil {
			data["status"] = "fail"
			data["msg"] = err.Error()
			jsonp(w, r, data)
			return
		}

		mod := model.GetModel(table)
		if mod == nil {
			log.Println("deleteTableByPKHandler: not found model " + table)
			data["status"] = "fail"
			data["msg"] = "no model"
			jsonp(w, r, data)
			return
		}
		_, err = mod.Delete(pk)
		if err != nil {
			data["status"] = "fail"
			data["msg"] = err.Error()
			jsonp(w, r, data)
			return
		}
		data["status"] = "ok"
		jsonp(w, r, data)
		return
	})
}

func MissHandler(w http.ResponseWriter, r *http.Request) {
	render.Render(w, "404.html", _locals(r, nil, true, true))
}

func checkRowsError(rows *sql.Rows) {
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
}

func nameNotin(f http.HandlerFunc, args ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			log.Println("input name empty")
			http.Error(w, "input name empty", http.StatusBadRequest)
			return
		}
		for _, arg := range args {
			if name == arg {
				log.Println("input name should not be", arg)
				http.Error(w, "input name forbidden", http.StatusBadRequest)
				return
			}
		}

		f(w, r)
	}
}

func checkErr(err error, args ...string) {
	var hint string
	if len(args) < 1 {
		hint = "Err is: "
	}
	if len(args) >= 1 {
		hint = args[0]
	}
	if err != nil {
		log.Fatal(hint, err)
	}
}

func ifUnexistsCreateDir(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filename, 0755)
		}
	}
}

func Paginator(page, prepage int, nums int64) map[string]interface{} {

	var prevpage int
	var nextpage int

	totalpages := int(math.Ceil(float64(nums) / float64(prepage)))
	if page > totalpages {
		page = totalpages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalpages-5 && totalpages > 5: //last five pages
		start := totalpages - 5 + 1
		prevpage = page - 1
		nextpage = int(math.Min(float64(totalpages), float64(page+1)))
		pages = make([]int, 5)
		for i, _ := range pages {
			pages[i] = start + i
		}
	case page >= 3 && totalpages > 5:
		start := page - 3 + 1
		pages = make([]int, 5)
		prevpage = page - 3
		for i, _ := range pages {
			pages[i] = start + i
		}
		prevpage = page - 1
		nextpage = page + 1
	default:
		pages = make([]int, int(math.Min(5, float64(totalpages))))
		for i, _ := range pages {
			pages[i] = i + 1
		}
		prevpage = int(math.Max(float64(1), float64(page-1)))
		nextpage = page + 1
	}
	paginatorMap := make(map[string]interface{})
	if page == totalpages {
		nextpage = totalpages
	}
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalpages
	paginatorMap["prevpage"] = prevpage
	paginatorMap["nextpage"] = nextpage
	paginatorMap["currpage"] = page
	paginatorMap["firstpage"] = 1
	paginatorMap["lastpage"] = totalpages
	return paginatorMap
}

func queryDataTotal(sqlString string, args ...interface{}) (dataCount int64) {
	countRows, err := app.db.Query(sqlString, args...)
	checkErr(err)
	defer countRows.Close()
	for countRows.Next() {
		var count int64
		err := countRows.Scan(&count)
		checkErr(err)
		dataCount = count
	}
	return
}

func getArgument(r *http.Request, want string, wantDefault interface{}) interface{} {
	argument := r.FormValue(want)
	if argument == "" {
		if wantDefault == nil {
			return nil
		}
		return wantDefault
	}
	return argument
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
