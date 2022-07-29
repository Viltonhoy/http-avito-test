// custom unmarshal to check the sorting condition of the user's transaction history
package storage

import (
	"encoding/json"
	"errors"
)

var ErrBadOrderType = errors.New("wrong value of ordBy type")

func (j *OrdBy) UnmarshalJSON(v []byte) error {
	var s string

	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}

	switch s {
	case string(OrderByAmount):
		*j = OrderByAmount
	case string(OrderByDate):
		*j = OrderByDate
	default:
		return ErrBadOrderType
	}

	return nil
}
