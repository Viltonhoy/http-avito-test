package generatetable

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type rowSrcPosting struct {
	rows []Posting
	idx  int
	err  error
}

func AddGenerateTable() int64 {

	rowSrcData := rowSrcPosting{
		rows: GenerateTableData(5, 100),
		idx:  -1,
		err:  nil,
	}

	columnName := []string{"account_id", "cb_journal", "accounting_period", "amount", "date", "addressee"}

	var tx pgx.Tx

	r, err := tx.CopyFrom(context.Background(), pgx.Identifier{"posting"}, columnName, &rowSrcData)
	if err != nil {
		return 0
	}
	return r
}

func (b *rowSrcPosting) Next() bool {
	b.idx++
	return b.idx < len(b.rows)
}

func (p Posting) interfaceSlice() ([]interface{}, error) {
	return []interface{}{
		p.accountID,
		p.CBjournal,
		p.accountingPeriod,
		p.amount,
		p.date,
		p.addressee,
	}, nil
}

func (b *rowSrcPosting) Values() ([]interface{}, error) {
	data, err := b.rows[b.idx].interfaceSlice()
	b.err = err
	return data, err
}

func (b *rowSrcPosting) Err() error {
	return b.err
}
