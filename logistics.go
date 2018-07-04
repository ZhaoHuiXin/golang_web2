package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.autoforce.net/autoforce/admin/model"
)

func init() {
	app.HandleFunc("GET", "/v1/ixiao/logistics/bill",
		logisticsManagerHandler)
	app.HandleFunc("GET", "/v1/ixiao/logistics/bill/detail",
		logisticsDetailHandler)
}

func logisticsManagerHandler(w http.ResponseWriter, r *http.Request) {
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
	begin = begin * 30

	dataCount := queryDataTotal("SELECT count(*) " +
		"from cheyixiao.logistics_bill")

	lgRows, err := app.db.Query("SELECT id,created_at,saler_id,"+
		"city_id_begin, city_id_end,send_time,car_type,car_price,"+
		"car_num FROM cheyixiao.logistics_bill "+
		"order by created_at desc LIMIT ?, 30", begin)

	checkErr(err)
	defer lgRows.Close()
	var lgs []model.LogisticsBill
	for lgRows.Next() {
		var lgId int
		var salerId sql.NullInt64
		var createdAt time.Time
		var cityIdBegin sql.NullInt64
		var cityIdEnd sql.NullInt64
		var sendTime time.Time
		var carType sql.NullInt64
		var carPrice sql.NullInt64
		var carNum sql.NullInt64
		err := lgRows.Scan(&lgId, &createdAt, &salerId, &cityIdBegin,
			&cityIdEnd, &sendTime, &carType, &carPrice, &carNum)
		checkErr(err)
		cityFromId := cityIdBegin.Int64
		cityToId := cityIdEnd.Int64
		cityFromRes, err := app.Query("sELECT name FROM "+
			"cheyixiao.areas WHERE id=?", cityFromId)
		checkErr(err)
		cityToRes, err := app.Query("sELECT name FROM "+
			"cheyixiao.areas WHERE id=?", cityToId)
		checkErr(err)
		dealer := ""
		if salerId.Int64 != 0 {
			dealerRes, err := app.Query("select cheyixiao.dealers.name "+
				"from cheyixiao.dealers join cheyixiao.salers on "+
				"cheyixiao.dealers.id=cheyixiao.salers.dealer_id "+
				"where cheyixiao.salers.id=?", salerId)
			checkErr(err)
			dealer = dealerRes[0]["name"].(string)
		}
		lg := model.LogisticsBill{
			Id:          lgId,
			SalerId:     salerId.Int64,
			CreatedAt:   createdAt.Format("2006-01-02 15:04:05"),
			CityIdBegin: cityFromRes[0]["name"].(string),
			CityIdEnd:   cityToRes[0]["name"].(string),
			SendTime:    sendTime.Format("2006-01-02 15:04:05"),
			CarType:     int(carType.Int64),
			CarPrice:    int(carPrice.Int64),
			CarNum:      int(carNum.Int64),
			DealerName:  dealer,
		}
		lgs = append(lgs, lg)
	}
	checkRowsError(lgRows)
	data["lgs"] = lgs
	res := Paginator(page, 30, dataCount)
	data["paginator"] = res
	data["start"] = start
	_locals(r, data, true, true).Render(w, "logistics.html")

}

func logisticsDetailHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	billId := r.FormValue("id")
	if billId == "" {
		data["status"] = "fail"
		data["msg"] = "lose necessary parameter id!"
		jsonp(w, r, data)
		return
	}
	querySql := fmt.Sprintf("SELECT id,created_at,"+
		"city_id_begin, city_id_end,send_time,car_type,car_price,"+
		"car_num,sender_name,sender_phone,sender_address,"+
		"receiver_name,receiver_phone,receiver_address,"+
		"is_invoice FROM cheyixiao.logistics_bill WHERE id=%s", billId)
	billRes, err := app.Query(querySql)
	billObj := billRes[0]
	checkErr(err)
	cityFromRes, err := app.Query("sELECT name FROM "+
		"cheyixiao.areas WHERE id=?", billObj["city_id_begin"].(string))
	checkErr(err)
	cityToRes, err := app.Query("sELECT name FROM "+
		"cheyixiao.areas WHERE id=?", billObj["city_id_end"].(string))
	checkErr(err)
	lgId, _ := strconv.Atoi(billObj["id"].(string))
	carType, _ := strconv.Atoi(billObj["car_type"].(string))
	carPrice, _ := strconv.Atoi(billObj["car_price"].(string))
	carNum, _ := strconv.Atoi(billObj["car_num"].(string))
	isInvoice, _ := strconv.Atoi(billObj["is_invoice"].(string))
	senderPhone, _ := strconv.ParseInt(billObj["sender_phone"].(string), 10, 64)
	receiverPhone, _ := strconv.ParseInt(billObj["receiver_phone"].(string), 10, 64)
	sendTime, _ := billObj["send_time"].(string)
	createdAt, _ := billObj["created_at"].(string)
	lg := model.LogisticsBill{
		Id:              lgId,
		CreatedAt:       createdAt,
		CityIdBegin:     cityFromRes[0]["name"].(string),
		CityIdEnd:       cityToRes[0]["name"].(string),
		SendTime:        sendTime,
		CarType:         carType,
		CarPrice:        carPrice,
		CarNum:          carNum,
		SenderName:      billObj["sender_name"].(string),
		SenderAddress:   billObj["sender_address"].(string),
		SenderPhone:     senderPhone,
		ReceiverName:    billObj["receiver_name"].(string),
		ReceiverAddress: billObj["receiver_address"].(string),
		ReceiverPhone:   receiverPhone,
		IsInvoice:       isInvoice,
	}
	data["lg"] = lg
	_locals(r, data, true, true).Render(w, "logistics_info.html")
}
