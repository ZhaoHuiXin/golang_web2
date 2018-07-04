package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"fmt"
	"time"
)

func init() {
	app.HandleFunc("GET", "/v1/ixiao/cars/quote",
		quotedPriceHandler)
	app.HandleFunc("PUT", "/v1/ixiao/cars/quote",
		modifyQuotedPriceHandler)
	app.HandleFunc("GET", "/v1/ixiao/cars/brands",
		carBrandsHandler)
	app.HandleFunc("GET", "/v1/ixiao/cars/series",
		carSeriesHandler)
	app.HandleFunc("GET", "/v1/ixiao/cars/specs",
		carSpecsHandler)
}

type Quote struct {
	SpecId         int
	BrandName       string
	SeriesName    string
	SpecName    string
	GuidePrice   string
	RealPrice      string
	UpdatedAt     string
}

func quotedPriceHandler(w http.ResponseWriter, r *http.Request) {
	var dataCount int64
	var rows *sql.Rows
	var table string
	var id interface{}
	var err error
	data := make(map[string]interface{})
	start := r.FormValue("start")
	var page int
	var begin int
	if start == "" {
		begin = 0
		page = 1
	} else {
		page, _ = strconv.Atoi(start)
		begin = page - 1
	}
	begin = begin * 10
	brandId := getArgument(r,"brandId", nil)
	seriesId := getArgument(r,"seriesId", nil)
	specId := getArgument(r,"specId", nil)
	data["brandId"] = "0"
	data["seriesId"] = "0"
	data["specId"] = "0"
	if brandId == nil && seriesId==nil && specId == nil{
		data["quotes"] = ""
		data["model"] = 0
		_locals(r, data, true, true).Render(w, "quoted_price.html")
		return
	}
	switch{
		case brandId!= nil:
			table = "t3"
			id = brandId
			data["brandId"] = brandId
		case seriesId!=nil:
			table = "t2"
			id = seriesId
			data["seriesId"] = seriesId
		case specId!=nil:
			table = "t1"
			id = specId
			data["specId"] = specId
	}
	countSql := fmt.Sprintf("select count(*) from " +
		"cheyixiao.specs as t1 join cheyixiao.series as t2 " +
		"on t1.series_id=t2.id join cheyixiao.brands as t3 " +
		"on t2.brand_id = t3.id where t3.status=1 and %s.id=?",table)
	dataCount = queryDataTotal(countSql, id)
	querySql := fmt.Sprintf("select t1.id,t3.name,t2.name,t1.name," +
		"t1.guide_price,t1.real_price,t1.updated_at from cheyixiao.specs as t1 " +
		"join cheyixiao.series as t2 on t1.series_id=t2.id join cheyixiao.brands " +
		"as t3 on t2.brand_id = t3.id where t3.status=1 and %s.id=? LIMIT ?,10", table)
	rows, err = app.db.Query(querySql, id, begin)
	checkErr(err)
	defer rows.Close()
	var quotes []Quote
	for rows.Next() {
		var specId  int
		var brandName      sql.NullString
		var seriesName    sql.NullString
		var specName    sql.NullString
		var guidePrice   sql.NullString
		var realPrice      sql.NullString
		var updatedAt     time.Time
		err := rows.Scan(&specId, &brandName, &seriesName, &specName, &guidePrice,
			&realPrice, &updatedAt)
		checkErr(err)
		quote := Quote{
			SpecId:  specId,
			BrandName:	brandName.String,
			SeriesName:    seriesName.String,
			SpecName:   specName.String,
			GuidePrice :  guidePrice.String,
			RealPrice :     realPrice.String,
			UpdatedAt :    updatedAt.Format("2006-01-02 15:04:05"),
		}
		quotes = append(quotes, quote)
	}
	checkRowsError(rows)
	data["quotes"] = quotes
	res := Paginator(page, 10, dataCount)
	data["paginator"] = res
	if brandId == nil && seriesId==nil && specId == nil{
		data["quotes"] = ""
	}
	data["model"] = 1
	_locals(r, data, true, true).Render(w, "quoted_price.html")
}

func modifyQuotedPriceHandler(w http.ResponseWriter, r *http.Request)  {
	data := make(map[string]interface{})
	specId := getArgument(r, "specId",nil)
	if specId == nil{
		data["status"] = "fail"
		data["msg"] = "缺失必要参数车型id"
		jsonp(w, r, data)
		return
	}
	realPrice := getArgument(r, "realPrice", "")
	specsToUpdate := make(map[string]map[string]interface{})
	specsAlternations := make(map[string]interface{})
	specsConditions := make(map[string]interface{})
	specsAlternations["real_price"] = realPrice
	specsConditions["id"] = specId
	specsToUpdate["$set"] = specsAlternations
	_, err :=app.FindOneAndUpdate("cheyixiao.specs", specsConditions, specsToUpdate, false)
	checkErr(err)
	data["status"] = "ok"
	jsonp(w, r, data)
}

func carBrandsHandler(w http.ResponseWriter, r *http.Request){
	data := make(map[string]interface{})
	res, err := app.Query("select id,name from cheyixiao.brands where status=1")
	checkErr(err)
	data["brands"] = res
	jsonp(w, r, data)

}

func carSeriesHandler(w http.ResponseWriter, r *http.Request){
	data := make(map[string]interface{})
	brandId := getArgument(r, "brandId", nil)
	res, err := app.Query("select id,name from " +
		"cheyixiao.series where brand_id=?", brandId)
	checkErr(err)
	data["series"] = res
	jsonp(w, r, data)
}

func carSpecsHandler(w http.ResponseWriter, r *http.Request){
	data := make(map[string]interface{})
	seriesId := getArgument(r, "seriesId", nil)
	res, err := app.Query("select id,name from cheyixiao.specs " +
		"where series_id=?", seriesId)
	checkErr(err)
	data["specs"] = res
	jsonp(w, r, data)
}
