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
	ID       int64     `json:"ID"`
	Type     string    `json:"Type"`
	Sum      int64     `json:"Sum"`
	Location string    `json:"Location"`
	Date     time.Time `json:"Date"`
}

type TransfList []Transf
