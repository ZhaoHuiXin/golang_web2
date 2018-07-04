package model

import (
	"database/sql"
)

type LogisticsBill struct {
	_Base
	Id              int
	SalerId         int64
	CityIdBegin     string
	CityIdEnd       string
	SendTime        string
	CarType         int
	CarPrice        int
	CarNum          int
	CreatedAt       string
	ReceiverName    string
	ReceiverPhone   int64
	ReceiverAddress string
	SenderName      string
	SenderPhone     int64
	SenderAddress   string
	IsInvoice       int
	DealerName      string
	TaxNumber       string
	ReceiverCompany string
}

var logistics_bill_config = &Config{
	db:    "cheyixiao",
	table: "logistics_bill",
	cols:  []columndef{},
}

func NewLogisticsBill(db *sql.DB) Model {
	p := &LogisticsBill{}
	p._Base.db = db
	p._Base.config = logistics_bill_config
	return p
}

func init() {
	Register("cheyixiao.logistics_bill", NewLogisticsBill)
}
