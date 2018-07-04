package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"git.autoforce.net/autoforce/admin/model"
)

func init() {
	app.HandleFunc("GET", "/v1/ixiao/finance/bill",
		billLoansHandler)
	app.HandleFunc("GET", "/v1/ixiao/finance/stock",
		stockLoansHandler)
}

func billLoansHandler(w http.ResponseWriter, r *http.Request) {
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
		"from cheyixiao.bill_loans")

	blRows, err := app.db.Query("SELECT id,saler_id,info,"+
		"short_name, car_type,car_id,car_color,car_num,car_price,"+
		"contract_fee,payed,pay_in,ratio,deposit,pic_purchase,pic_payed,"+
		"pic_procedure,created_at FROM cheyixiao.bill_loans LIMIT ?,30", begin)
	if err != nil {
		log.Fatal(err)
	}
	defer blRows.Close()
	var bills []model.BillLoans

	for blRows.Next() {
		var id int
		var salerId sql.NullInt64
		var info sql.NullString
		var shortName sql.NullString
		var carType int
		var carId int
		var carColor sql.NullString
		var carNum int
		var carPrice int
		var contractFee int
		var payed int
		var payIn int
		var ratio int
		var deposit int
		var picPurchase sql.NullString
		var picPayed sql.NullString
		var picProcedure sql.NullString
		var createdAt time.Time

		err := blRows.Scan(&id, &salerId, &info, &shortName, &carType,
			&carId, &carColor, &carNum, &carPrice, &contractFee, &payed,
			&payIn, &ratio, &deposit, &picPurchase, &picPayed, &picProcedure,
			&createdAt)
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
		res, err := app.Query("select i1,i2 from afsaas.specs where id=?", carId)
		carName := res[0]["i1"].(string)
		carGuide := res[0]["i2"].(string)
		bill := model.BillLoans{
			Id:           id,
			DealerName:   dealer,
			Info:         info.String,
			ShortName:    shortName.String,
			CreatedAt:    createdAt.Format("2006-01-02 15:04:05"),
			CarType:      carType,
			CarId:        carId,
			CarColor:     carColor.String,
			CarNum:       carNum,
			CarPrice:     carPrice,
			ContractFee:  contractFee,
			Payed:        payed,
			PayIn:        payIn,
			Ratio:        ratio,
			Deposit:      deposit,
			PicPurchase:  picPurchase.String,
			PicPayed:     picPayed.String,
			PicProcedure: picProcedure.String,
			CarName:      carName,
			CarGuide:     carGuide,
		}
		bills = append(bills, bill)
	}
	checkRowsError(blRows)
	data["bills"] = bills
	res := Paginator(page, 30, dataCount)
	data["paginator"] = res
	data["type"] = "bill"
	_locals(r, data, true, true).Render(w, "finance.html")
}

func stockLoansHandler(w http.ResponseWriter, r *http.Request) {
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
		"from cheyixiao.stock_loans")
	slRows, err := app.db.Query("SELECT id,saler_id,info,"+
		"total_fee, pay_in,ratio,deposit,pic_cert,pic_invoice,"+
		"created_at FROM cheyixiao.stock_loans "+
		"LIMIT ?,30", begin)
	checkErr(err)
	defer slRows.Close()
	var bills []model.StockLoans
	for slRows.Next() {
		var id int
		var salerId sql.NullInt64
		var info sql.NullString
		var totalFee int
		var payIn int
		var ratio int
		var deposit int
		var picCert sql.NullString
		var picInvoice sql.NullString
		var createdAt time.Time
		err := slRows.Scan(&id, &salerId, &info, &totalFee, &payIn, &ratio,
			&deposit, &picCert, &picInvoice, &createdAt)
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
		cars := stockLoansCarsHandler(id)
		bill := model.StockLoans{
			Id:         id,
			DealerName: dealer,
			Info:       info.String,
			CreatedAt:  createdAt.Format("2006-01-02 15:04:05"),
			PayIn:      payIn,
			Ratio:      ratio,
			Deposit:    deposit,
			TotalFee:   totalFee,
			PicCert:    picCert.String,
			PicInvoice: picInvoice.String,
			Cars:       cars,
		}
		bills = append(bills, bill)
	}
	checkRowsError(slRows)
	data["bills"] = bills
	data["type"] = "stock"

	res := Paginator(page, 30, dataCount)
	data["paginator"] = res
	_locals(r, data, true, true).Render(w, "finance.html")
}

func stockLoansCarsHandler(id int) []map[string]interface{} {
	res, err := app.Query("select t1.id,car_id,t3.i1, t3.i2, car_type,"+
		"car_color,car_price,car_num from cheyixiao.stock_loans_cars as t1 "+
		"join cheyixiao.stock_loans as t2 on t1.sl_id=t2.id "+
		"join afsaas.specs t3 on t1.car_id=t3.id where t2.id=?", id)
	checkErr(err)
	return res
}
