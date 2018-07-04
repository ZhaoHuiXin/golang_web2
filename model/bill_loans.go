package model

import (
	"database/sql"
)

type BillLoans struct {
	_Base
	Id           int
	SalerId      int64
	Info         string
	ShortName    string
	CarType      int
	CarId        int
	CarColor     string
	CarNum       int
	CarPrice     int
	ContractFee  int
	Payed        int
	PayIn        int
	Ratio        int
	Deposit      int
	PicPurchase  string
	PicPayed     string
	PicProcedure string
	CreatedAt    string
	DealerName   string
	CarName      string
	CarGuide     string
}

var bill_loans_config = &Config{
	db:    "cheyixiao",
	table: "bill_loans",
	cols:  []columndef{},
}

func NewBillLoans(db *sql.DB) Model {
	p := &BillLoans{}
	p._Base.db = db
	p._Base.config = bill_loans_config
	return p
}

func init() {
	Register("cheyixiao.bill_loans", NewBillLoans)
}
