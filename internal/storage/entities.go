package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID      int64
	Balance decimal.Decimal
}

type Transf struct {
	ID       int64
	Type     string
	Sum      int64
	Location string
	Date     time.Time
}

type TransfList struct {
	List []Transf
}
