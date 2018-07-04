package model

import (
	"database/sql"
)

type StockLoans struct {
	_Base
	Id         int
	SalerId    int64
	Info       string
	ShortName  string
	TotalFee   int
	PayIn      int
	Ratio      int
	Deposit    int
	PicCert    string
	PicInvoice string
	CreatedAt  string
	DealerName string
	Cars       []map[string]interface{}
}

var stock_loans_config = &Config{
	db:    "cheyixiao",
	table: "stock_loans",
	cols:  []columndef{},
}

func NewStockLoans(db *sql.DB) Model {
	p := &StockLoans{}
	p._Base.db = db
	p._Base.config = stock_loans_config
	return p
}

func init() {
	Register("cheyixiao.stock_loans", NewStockLoans)
}
