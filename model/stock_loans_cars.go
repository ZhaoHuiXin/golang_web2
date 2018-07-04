package model

import (
	"database/sql"
)

type StockLoansCars struct {
	_Base
	Id       int
	SlId     int
	CarType  int
	CarId    int
	CarColor string
	CarNum   int
	CarPrice int
	CarName  string
	CarGuide string
}

var stock_loans_cars_config = &Config{
	db:    "cheyixiao",
	table: "stock_loans_cars",
	cols:  []columndef{},
}

func NewStockLoansCars(db *sql.DB) Model {
	p := &StockLoansCars{}
	p._Base.db = db
	p._Base.config = stock_loans_cars_config
	return p
}

func init() {
	Register("cheyixiao.stock_loans_cars", NewStockLoansCars)
}
