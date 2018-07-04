package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"reflect"
	"time"

	"go/types"

	"git.autoforce.net/autoforce/admin/model"
	"github.com/RoaringBitmap/roaring"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	db_dsn string
	db     *sql.DB
	router *mux.Router

	features   []*Feature
	roleAccess map[int]*roaring.Bitmap
	userAccess map[int]*roaring.Bitmap
}

type Feature struct {
	model.Feature
}

var app = NewApp()

func init() {
	flag.StringVar(&app.db_dsn, "mysql", "mysql://user:pass@tcp(host:port)/dbname", "mysql uri")
}

func NewApp() *App {
	return &App{
		router: mux.NewRouter(),
	}
}

func (p *App) Init() {
	p.Open()
	model.RegisterDB(p.db)
	p.FeatureSync()
	p.LoadAccess()
}

func (p *App) Open() (err error) {
	if p.db != nil {
		return
	}

	u, err := url.Parse(p.db_dsn)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("parseTime", "true")
	q.Set("loc", "Local")
	u.RawQuery = q.Encode()

	log.Println("conn db", u.String())
	db, err := sql.Open("mysql", u.String())
	if err != nil {
		log.Fatal("open db fail", err)
		return
	}
	p.db = db
	db.SetMaxOpenConns(100)
	return
}

func (p *App) FeatureSync() {
	for _, feat := range p.features {
		feat.Sync()
	}
}

func (p *App) GetFeatureIndex(methods, path string) int {
	for i, feat := range p.features {
		if feat.Methods == methods && feat.Path == path {
			return i
		}
	}
	return -1
}

func (p *App) LoadAccess() {
	feat_id_index := make(map[int]int)
	roleAccess := make(map[int]*roaring.Bitmap)
	userAccess := make(map[int]*roaring.Bitmap)

	{
		rows, err := p.db.Query("SELECT `id`, `methods`, `path` FROM features")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var methods string
			var path string
			if err := rows.Scan(&id, &methods, &path); err != nil {
				log.Println(err)
				continue
			}
			index := p.GetFeatureIndex(methods, path)
			if index != -1 {
				feat_id_index[id] = index
			}
		}
		checkRowsError(rows)
	}

	{
		rows, err := p.db.Query("SELECT role_id,user_id,feat_id FROM access")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var feat_id int
			var role_id int
			var user_id int
			if err := rows.Scan(&role_id, &user_id, &feat_id); err != nil {
				log.Println(err)
			}
			var feat_index int
			if index, has := feat_id_index[feat_id]; has {
				feat_index = index
			} else {
				continue
			}
			if role_id > 0 {
				if _, has := roleAccess[role_id]; !has {
					roleAccess[role_id] = roaring.New()
				}
				roleAccess[role_id].Add(uint32(feat_index))
			}
			if user_id > 0 {
				if _, has := userAccess[user_id]; !has {
					userAccess[user_id] = roaring.New()
				}
				userAccess[user_id].Add(uint32(feat_index))
			}
		}
		checkRowsError(rows)
		p.roleAccess = roleAccess
		p.userAccess = userAccess
	}
}

// use feat index of features, for less mem
func (p *App) AuthorizedAsRole(role_id, feat_idx int) bool {
	if access, ok := p.roleAccess[role_id]; ok && access != nil {
		return access.Contains(uint32(feat_idx))
	}
	return false
}

func (p *App) AuthorizedAsUser(uid, feat_idx int) bool {
	if access, ok := p.userAccess[uid]; ok && access != nil {
		return access.Contains(uint32(feat_idx))
	}
	return false
}

func (p *App) AccessHandler(feat *Feature, f http.HandlerFunc) http.HandlerFunc {
	index := p.GetFeatureIndex(feat.Methods, feat.Path)
	return func(w http.ResponseWriter, r *http.Request) {
		if !feat.Auth {
			f(w, r)
			return
		}
		session, err := getSession(r)
		if err != nil || session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		var uid int
		var role_id int
		if v := session.Values["uid"]; v != nil {
			uid = v.(int)
		}
		if v := session.Values["role_id"]; v != nil {
			role_id = v.(int)
		}
		if uid < 1 {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		if p.AuthorizedAsUser(uid, index) || p.AuthorizedAsRole(role_id, index) {
			f(w, r)
			return
		}
		if opt.debug {
			f(w, r)
			return
		}
		render.Render(w, "forbidden.html", nil)
	}
}

func (p *App) HandleFunc(methods, path string, f http.HandlerFunc) {
	feat := &Feature{
		Feature: model.Feature{
			Path:    path,
			Methods: methods,
			Auth:    true,
		},
	}
	p.features = append(p.features, feat)
	p.router.HandleFunc(feat.Path, p.AccessHandler(feat, f)).Methods(feat.Methods)
}

func (p *App) Close() {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
}

func (p *App) Count(table, where string, args ...interface{}) int {
	num := 0
	s := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s", table, where)
	err := p.db.QueryRow(s, args...).Scan(&num)
	if err != nil {
		log.Println("QueryRow", err, s, args)
	}
	return num
}

func (p *App) Query(sqlString string, args ...interface{}) (res []map[string]interface{}, err error) {
	rows, err := p.db.Query(sqlString, args...)
	if err != nil {
		return
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for i := 0; i < count; i++ {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if val != nil {
				dataType := reflect.TypeOf(val).String()
				if dataType == "time.Time" {
					v = val.(time.Time).Format("2006-01-02 15:04:05")
					entry[col] = v
					continue
				}
			}
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		res = append(res, entry)
	}

	if err := rows.Err(); err != nil {
		log.Println("app query", err)
		err = nil
	}
	return
}

func (p *App) FindOneAndUpdate(tableName string, querys map[string]interface{}, args map[string]map[string]interface{}, upsert bool) (res []map[string]interface{}, err error) {
	var queryConditions string
	var queryParameterSlice []interface{}

	queryCount := 1
	for k, v := range querys {
		if queryCount == 1 {
			queryConditions += ("`" + k + "`=?")
			queryParameterSlice = append(queryParameterSlice, v)
			queryCount += 1
			continue
		}
		queryConditions += (",`" + k + "`=?")
		queryParameterSlice = append(queryParameterSlice, v)
	}
	querySql := "select * from " + tableName + " where " + queryConditions
	rowExists, err := p.Query(querySql, queryParameterSlice...)
	checkErr(err)

	if rowExists == nil {
		fmt.Println("upsert insert model")
		var info map[string]interface{}
		if upsert == true {
			fmt.Println("TRUE")
			fmt.Println(args)
			if _, ok := args["$set"]; ok {
				fmt.Println("$SET")
				info = args["$set"]
				if _, ok := info["id"]; ok {
					if info["id"] == nil {
						delete(info, "id")
					}
				}
				for k, v := range querys {
					switch v.(type) {
					case types.Map:
						continue
					case types.Array:
						continue
					case types.Tuple:
						continue
					case types.Slice:
						continue
					}
					if k == "id" && querys[k] != nil {
						continue
					}
					info[k] = v
				}
				fmt.Println(info)
				res, err = p.InsertOne(tableName, info)
				checkErr(err)
				fmt.Println(res)
				return
			}
			fmt.Println("nonono")
			return rowExists, err
		}
	}
	if _, ok := args["$set"]; ok {
		fmt.Println("update model")
		var conditions string
		var alterations string
		var updateParameterSlice []interface{}
		var backQueryParameterSlice []interface{}
		coditCount := 1
		alterCount := 1

		for k, v := range args["$set"] {
			if alterCount == 1 {
				alterations += ("`" + k + "`=?")
				updateParameterSlice = append(updateParameterSlice, v)
				alterCount += 1
				continue
			}
			alterations += (",`" + k + "`=?")
			updateParameterSlice = append(updateParameterSlice, v)
		}
		for k, v := range querys {
			if coditCount == 1 {
				conditions += ("`" + k + "`=?")
				updateParameterSlice = append(updateParameterSlice, v)
				backQueryParameterSlice = append(backQueryParameterSlice, v)
				coditCount += 1
				continue
			}
			conditions += (",`" + k + "`=?")
			updateParameterSlice = append(updateParameterSlice, v)
			backQueryParameterSlice = append(backQueryParameterSlice, v)
		}
		updateSql := "update " + tableName + " set " + alterations + " where " + conditions + " limit 1"
		fmt.Println(updateSql)

		stmt, err := p.db.Prepare(updateSql)
		checkErr(err, "FindOneAndUpdate define stmt error: ")
		rdata, err := stmt.Exec(updateParameterSlice...)
		fmt.Println(rdata.LastInsertId())

		checkErr(err, "FindOneAndUpdate stmt.Exec error: ")
		backQuerySql := "select * from " + tableName + " where " + conditions
		res, err = p.Query(backQuerySql, backQueryParameterSlice...)
		checkErr(err)
		return res, err
	}
	return
}

func (p *App) InsertOne(tableName string, args map[string]interface{}) (res []map[string]interface{}, err error) {
	var columns string
	var placeholder string
	var parameterSlice []interface{}
	count := 1
	for k, v := range args {
		if count == 1 {
			columns += ("`" + k + "`")
			placeholder += "?"
			parameterSlice = append(parameterSlice, v)
			count += 1
			continue
		}
		columns += (",`" + k + "`")
		placeholder += ",?"
		parameterSlice = append(parameterSlice, v)
	}
	insertSql := "insert into " + tableName + "(" + columns + ")" + " values " +
		"(" + placeholder + ")"
	fmt.Println(insertSql)
	stmt, err := p.db.Prepare(insertSql)
	checkErr(err, "InsertOne define stmt error: ")
	fmt.Println(parameterSlice...)
	rdata, err := stmt.Exec(parameterSlice...)
	id, _ := rdata.LastInsertId()
	fmt.Println(id)
	checkErr(err, "InsertOne stmt.Exec error: ")
	//fmt.Println(parameterSlice)
	querySql := "select * from " + tableName + " where id=?"
	res, err = p.Query(querySql, id)
	checkErr(err)
	return
}
