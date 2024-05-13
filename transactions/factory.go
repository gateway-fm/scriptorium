package transactions

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type PgTransactionFactory struct {
	*pg.DB
}

func NewPgTransactionFactory(db *pg.DB) *PgTransactionFactory {
	return &PgTransactionFactory{db}
}

func (f *PgTransactionFactory) Transaction(ctx context.Context) Transaction {
	tx, ok := ctx.Value(txKey).(Transaction)
	if !ok {
		return f.DB
	}

	return tx
}
