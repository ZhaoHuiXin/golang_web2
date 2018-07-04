package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	app.HandleFunc("GET", "/v1/ixiao/statistics/manager",
		statisticsManagerHandler)
}
func statisticsManagerHandler(w http.ResponseWriter, r *http.Request) {
	var querySql string
	var q, table string
	today := time.Now().Format("2006-01-02")
	s := strings.Split(today, "-")
	firstDay := s[0] + "-" + s[1] + "-" + "01"
	model := getArgument(r, "model", "register")
	begin := getArgument(r, "begin", firstDay)
	end := getArgument(r, "end", today)
	endStr := end.(string)
	endStr += " 23:59:59"
	switch model {
	case "active", "activeAddUp":
		q = "count(distinct(user_id))"
		table = "cheyixiao.look"
	case "register", "registerAddUp":
		q = "count(id)"
		table = "cheyixiao.salers"
	case "deal", "dealAddUp":
		q = "count(id)"
		table = "cheyixiao.user_status"
	case "dealSaler", "dealSalerAddUp":
		q = "count(distinct(saler_id))"
		table = "cheyixiao.user_status"
	}
	querySql = fmt.Sprintf("select %s as count,"+
		"date_format(created_at,'%%Y/%%m/%%d') as tm from "+
		"%s where created_at between '%s' and '%s' "+
		"group by tm", q, table, begin, endStr)
	//fmt.Println(querySql)
	data := make(map[string]interface{})
	res, err := app.Query(querySql)
	checkErr(err, "QUERY ERROR: ")

	if model == "register" || model == "active" || model == "deal" ||
		model == "dealSaler" {
		data["status"] = "ok"
		data["data"] = res
		_locals(r, data, true, true).Render(w, "statistics.html")
		return
	}
	addUpQuerySql := fmt.Sprintf("SELECT %s as count FROM "+
		"%s where created_at < '%s'", q, table, begin)
	resAddUp, err := app.Query(addUpQuerySql)
	checkErr(err, "QUERY ERROR: ")
	basicNum, _ := strconv.Atoi(resAddUp[0]["count"].(string))
	//fmt.Println(basicNum, reflect.TypeOf(basicNum))
	for _, v := range res {
		addNum, _ := strconv.Atoi(v["count"].(string))
		basicNum += addNum
		v["count"] = basicNum
	}
	data["status"] = "ok"
	data["data"] = res
	//fmt.Println(data["data"])
	_locals(r, data, true, true).Render(w, "statistics.html")
}
