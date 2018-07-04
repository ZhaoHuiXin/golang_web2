package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"git.autoforce.net/autoforce/admin/model"
	"io/ioutil"
	"encoding/json"
)

func init() {
	app.HandleFunc("GET", "/v1/ixiao/register/manager",
		registerManagerHandler)
	app.HandleFunc("GET", "/v1/ixiao/register/info",
		registerInfoManagerHandler)
	app.HandleFunc("POST", "/v1/ixiao/register/info",
		updateDealerInfoManagerHandler)
	app.HandleFunc("GET", "/v1/ixiao/register/info/area",
		areaHandler)
	app.HandleFunc("PUT", "/v1/ixiao/register/manager/pass",
		updateTableByPKNoLoginHandler(`cheyixiao.salers`, `cert_code`))
	app.HandleFunc("PUT", "/v1/ixiao/register/manager/back",
		updateTableByPKNoLoginHandler(`cheyixiao.salers`, `reason`))
	app.HandleFunc("POST", "/v1/ixiao/img/upload",
		uploadImgHandler)
}

func registerManagerHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	start := r.FormValue("start")
	//fmt.Println("start: ", start)
	var page int
	var begin int
	if start == "" {
		begin = 0
		page = 1
	} else {
		page, _ = strconv.Atoi(start)
		begin = page - 1
	}
	begin = begin * 30
	//fmt.Println("begin: ", begin)
	//fmt.Println("page: ", page)
	cert := r.FormValue("cert")
	if cert == "" {
		cert = "1"
	}
	data["model"] = 0
	if cert == "1" {
		data["model"] = 1
	}

	dataCount := queryDataTotal("SELECT count(*) from cheyixiao.salers "+
		"WHERE cert_code=?", cert)

	rows, err := app.db.Query("SELECT id,username,role,updated_at,brands "+
		"FROM cheyixiao.salers WHERE cert_code=? "+
		"LIMIT ?, 30", cert, begin)
	// order by updated_at desc
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var salers []model.Saler
	for rows.Next() {
		var id int
		var username sql.NullString
		var brands sql.NullString
		var updated_at time.Time

		var role sql.NullInt64
		if err := rows.Scan(&id, &username, &role, &updated_at, &brands); err != nil {
			log.Fatal(err)
		}

		// add brands info to reponse data,
		// attention : don't use Query("select ..") directly, -_-!
		choose_brands := ""
		var selected_brands string
		if selected_brands = brands.String; selected_brands != "" {
			query_sql := fmt.Sprintf("SELECT name FROM "+
				"cheyixiao.brands WHERE id in(%s) ", selected_brands)
			brows, berr := app.db.Query(query_sql)
			if berr != nil {
				log.Fatal(berr)
			}
			defer brows.Close()
			for brows.Next() {
				var name string
				if err := brows.Scan(&name); err != nil {
					log.Fatal(err)
				}
				choose_brands += (name + " ")
			}
		}
		saler := model.Saler{
			Id:        id,
			Username:  username.String,
			Role:      role.Int64,
			UpdatedAt: updated_at.Format("2006-01-02 15:04:05"),
			Brands:    choose_brands,
		}
		salers = append(salers, saler)
	}
	checkRowsError(rows)
	data["salers"] = salers
	res := Paginator(page, 30, dataCount)
	data["paginator"] = res
	data["cert"] = cert
	data["start"] = start
	_locals(r, data, true, true).Render(w, "register_manager.html")
}

func registerInfoManagerHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	saler_id := r.FormValue("saler_id")
	cert := r.FormValue("cert")
	start := r.FormValue("start")
	if saler_id == "" {
		data["code"] = 204
		data["msg"] = "lose arg saler_id"
		jsonp(w, r, data)
		return
	}
	rows, err := app.db.Query("SELECT id, name, avatar, address, phone,"+
		"role, ID_face, ID_con, cert_code, city, dealer_id "+
		"FROM cheyixiao.salers WHERE id=? ", saler_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name sql.NullString
		var avatar sql.NullString
		var address sql.NullString
		var phone sql.NullInt64
		var ID_face sql.NullString
		var ID_con sql.NullString
		var role sql.NullInt64
		var cert_code sql.NullInt64
		var city sql.NullString
		var dealer_id sql.NullInt64
		var cityName string

		if err := rows.Scan(&id, &name, &avatar, &address, &phone, &role,
			&ID_face, &ID_con, &cert_code, &city, &dealer_id); err != nil {
			log.Fatal(err)
		}
		//fmt.Println(1, city.String, 1)
		salerRole := role.Int64
		dealerId := dealer_id.Int64
		if city.String != "" {
			res, err := app.Query("select t1.name as province,t2.name "+
				"as city from cheyixiao.areas as t1 join cheyixiao.areas "+
				"as t2 on t1.id=t2.pid where t2.id=?", city.String)
			checkErr(err)
			province := res[0]["province"].(string)
			city := res[0]["city"].(string)
			cityName = province + city
		}
		showPhone := phone.Int64
		phoneNumber, isReplace := readPhoneNumber(id)
		if isReplace == true{
			showPhone = phoneNumber
		}
		saler := model.Saler{
			Id:       id,
			Name:     name.String,
			Avatar:   avatar.String,
			Address:  address.String,
			Phone:    showPhone,
			IDFace:   ID_face.String,
			IDCon:    ID_con.String,
			Role:     salerRole,
			City:     city.String,
			CityName: cityName,
		}
		data["saler"] = saler
		//fmt.Println(dealerId)
		if salerRole == 1 || salerRole == 2 {
			dealer_rows, err := app.db.Query("SELECT id, `name`, company, "+
				"`call`, phone, ch_name, bus_licence, pic_door, pic_show, "+
				"pic_rest, pic_other, address "+
				"FROM cheyixiao.dealers WHERE id=?", dealerId)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println()
			for dealer_rows.Next() {
				var drId int
				var drName sql.NullString
				var company sql.NullString
				var drPhone sql.NullInt64
				var call sql.NullString
				var chName sql.NullString
				var busLicence sql.NullString
				var picDoor sql.NullString
				var picShow sql.NullString
				var picRest sql.NullString
				var picOther sql.NullString
				var dealerAddress sql.NullString
				if err := dealer_rows.Scan(&drId, &drName, &company, &call,
					&drPhone, &chName, &busLicence, &picDoor, &picShow, &picRest,
					&picOther, &dealerAddress); err != nil {
					log.Fatal(err)
				}
				showDrPhone := drPhone.Int64
				if isReplace == true{
					showDrPhone = phoneNumber
				}

				dealer := model.Dealer{
					Id:         drId,
					Name:       drName.String,
					Company:    company.String,
					Call:       call.String,
					Phone:      showDrPhone,
					ChName:     chName.String,
					BusLicence: busLicence.String,
					PicDoor:    picDoor.String,
					PicShow:    picShow.String,
					PicRest:    picRest.String,
					PicOther:   picOther.String,
					Address:    dealerAddress.String,
				}
				data["dealer"] = dealer
			}
		}
	}
	data["start"] = start
	data["cert"] = cert
	_locals(r, data, true, true).Render(w, "register_info.html")
}

func updateTableByPKNoLoginHandler(table, keyword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})
		k := r.FormValue("name")
		if k != keyword {
			data["status"] = "fail"
			data["msg"] = fmt.Sprintf("the keyword must be %s only !",
				keyword)
			jsonp(w, r, data)
			return
		}

		pk, err := strconv.ParseInt(r.FormValue("pk"), 10, 64)
		if err != nil {
			fmt.Println("cant convert", err)
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
		if keyword == "reason" {
			_, err = mod.Update(pk, sql.Named("cert_code", 3))
			if err != nil {
				data["status"] = "fail"
				data["msg"] = err.Error()
				jsonp(w, r, data)
				return
			}
		}
		data["status"] = "ok"
		jsonp(w, r, data)
		return
	}
}

func uploadImgHandler(w http.ResponseWriter, r *http.Request) {
	var uploadPath, picUrl string
	IXIAO_PROJECT := os.Getenv("IXIAO_PROJECT")
	IXIAO_DOMAIN := os.Getenv("IXIAO_DOMAIN")
	uploadPath = IXIAO_PROJECT + "/static/upload"
	picUrl = IXIAO_DOMAIN + "/static/upload"
	data := make(map[string]interface{})
	img_type := r.FormValue("type")
	switch img_type {
	case "avatar":
		uploadPath = uploadPath + "/avatar/"
		picUrl = picUrl + "/avatar/"
	case "IDFace", "IDCon":
		uploadPath = uploadPath + "/ID/"
		picUrl = picUrl + "/ID/"
	case "busLicence":
		uploadPath = uploadPath + "/license/"
		picUrl = picUrl + "/license/"
	case "picDoor", "picShow", "picRest", "picOther":
		uploadPath = uploadPath + "/pic/"
		picUrl = picUrl + "/pic/"
	}
	// create file path
	ifUnexistsCreateDir(uploadPath)
	file, head, err := r.FormFile("img")
	checkErr(err)
	defer file.Close()
	// generate md5ImgName
	suffix := path.Ext(head.Filename)
	timeNow := strconv.FormatInt(time.Now().Unix(), 10)
	md5ImgName := md5str(head.Filename+timeNow) + suffix
	fileUrl := picUrl + md5ImgName
	//fmt.Println("fileUrl: ", fileUrl)
	fW, err := os.Create(uploadPath + md5ImgName)
	if err != nil {
		fmt.Println("文件创建失败", err)
		return
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		fmt.Println("文件保存失败", err)
		return
	}
	//fileUrl := picUrl + md5ImgName
	//fmt.Println("fileUrl: ", fileUrl)
	data["status"] = "ok"
	data["imgUrl"] = fileUrl
	jsonp(w, r, data)
	return
}

func areaHandler(w http.ResponseWriter, r *http.Request) {
	pid := r.FormValue("pid")
	if pid == "" {
		pid = "0"
	}
	var data []map[string]int
	rows, err := app.db.Query("select name, id from cheyixiao.areas "+
		"where pid=? order by pinyin", pid)
	checkErr(err, "【areaHandler】Query from cheyixiao.areas Error: ")
	for rows.Next() {
		tmpMap := make(map[string]int)
		var name string
		var id int
		err := rows.Scan(&name, &id)
		checkErr(err, "【areaHandler】Scan cheyixiao.areas.rows Error: ")
		tmpMap[name] = id
		data = append(data, tmpMap)
	}
	jsonp(w, r, data)
}

func updateDealerInfoManagerHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	role := getArgument(r, "role", nil)
	if role == nil {
		data["status"] = "fail"
		data["msg"] = "lose necessary argument role"
		jsonp(w, r, data)
		return
	}
	salerId := getArgument(r, "salerId", nil)
	if salerId == nil {
		data["status"] = "fail"
		data["msg"] = "lose necessary argument salerId"
		jsonp(w, r, data)
		return
	}
	name := getArgument(r, "name", nil)
	contactAddress := getArgument(r, "contactAddress", nil)
	contactMethod := getArgument(r, "contactMethod", nil)
	avatar := getArgument(r, "avatar", nil)
	cityId := getArgument(r, "cityId", nil)

	var salersToUpdate map[string]map[string]interface{}
	salersToUpdate = make(map[string]map[string]interface{})
	var salersAlternations map[string]interface{}
	salersAlternations = make(map[string]interface{})
	var salersConditions map[string]interface{}
	salersConditions = make(map[string]interface{})
	salersAlternations["name"] = name
	salersAlternations["phone"] = contactMethod
	salersAlternations["address"] = contactAddress
	salersAlternations["avatar"] = avatar
	salersAlternations["city"] = cityId
	salersConditions["id"] = salerId
	salersToUpdate["$set"] = salersAlternations
	if role == "dealer" {
		app.FindOneAndUpdate("cheyixiao.salers", salersConditions, salersToUpdate, false)
		dealersToUpdate := make(map[string]map[string]interface{})
		dealersAlternations := make(map[string]interface{})
		dealersConditions := make(map[string]interface{})
		dealerId := getArgument(r, "dealerId", nil)
		if dealerId == "" {
			data["status"] = "fail"
			data["msg"] = "lose necessary argument dealerId"
			jsonp(w, r, data)
			return
		}
		dealersConditions["id"] = dealerId
		dealersAlternations["name"] = getArgument(r, "dealerName", nil)
		dealersAlternations["company"] = getArgument(r, "company", nil)
		dealersAlternations["address"] = getArgument(r, "dealerAddress", nil)
		dealersAlternations["call"] = getArgument(r, "dealerCall", nil)
		dealersAlternations["ch_name"] = getArgument(r, "chName", nil)
		dealersAlternations["phone"] = getArgument(r, "chPhone", nil)
		dealersAlternations["pic_door"] = getArgument(r, "picDoor", nil)
		dealersAlternations["pic_show"] = getArgument(r, "picShow", nil)
		dealersAlternations["pic_rest"] = getArgument(r, "picRest", nil)
		dealersAlternations["pic_other"] = getArgument(r, "picOther", nil)
		dealersAlternations["bus_licence"] = getArgument(r, "busLicence", nil)
		dealersToUpdate["$set"] = dealersAlternations
		app.FindOneAndUpdate("cheyixiao.dealers", dealersConditions, dealersToUpdate, false)
	}
	if role == "saler" {
		salersAlternations["ID_face"] = getArgument(r, "IDFace", nil)
		salersAlternations["ID_con"] = getArgument(r, "IDCon", nil)
		app.FindOneAndUpdate("cheyixiao.salers", salersConditions, salersToUpdate, false)
	}
	data["status"] = "ok"
	jsonp(w, r, data)
	return
}

type PhoneNumber struct {
	Remember		map[string]int64 `json:"remember"`
	Count      int `json:"count"`
	Numbers     []int64 `json:"numbers"`
}
func readPhoneNumber(salerId int)(pn int64, still bool){
	id := strconv.Itoa(salerId)
	isExist := checkFileIsExist("phoneNumber.json")
	if isExist == false{
		return 0,false
	}
	dataFormat := PhoneNumber{}
	cont, err := ioutil.ReadFile("phoneNumber.json")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(cont), &dataFormat)
	count := dataFormat.Count
	numbers := dataFormat.Numbers
	rememberMap := dataFormat.Remember
	if len(rememberMap) >0 {
		for k, v := range rememberMap{
			if k == id{
				return v, true
			}
		}
	}
	if count == len(numbers){
		return 0,false
	}
	pn = numbers[count]
	//fmt.Println("now: ", count, numbers[count])
	//fmt.Println("remember", dataFormat.Remember, reflect.TypeOf(dataFormat.Remember))
	rememberMap[id] = pn
	count += 1
	newData := PhoneNumber{
		Count:count,
		Numbers:numbers,
		Remember:rememberMap,
	}
	jsonData, err := json.Marshal(newData)
	checkErr(err)
	//fmt.Print([]byte(jsonData))

	os.Remove("phoneNumber.json")
	err = ioutil.WriteFile("phoneNumber.json", jsonData, 0666)
	checkErr(err)
	return pn, true
}

