package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID      int64           `json:"ID"`
	Balance decimal.Decimal `json:"Balance"`
}

type Transf struct {
	ID        int64     `json:"ID"`
	Type      string    `json:"Type"`
	Sum       int64     `json:"Sum"`
	Date      time.Time `json:"Date"`
	Addressee string    `json:"Addressee"`
}

type TransfList []Transf

type Balance struct {
	ID  int64
	Sum decimal.Decimal
}
