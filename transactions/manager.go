package transactions

import (
	"context"
)

type Options struct {
	AlwaysRollback bool
}

type PgTransactionManager struct {
	trf     *PgTransactionFactory
	options Options
}

func NewPgTransactionManager(trf *PgTransactionFactory, options Options) *PgTransactionManager {
	return &PgTransactionManager{
		trf:     trf,
		options: options,
	}
}

// Do function creates Tx for the given tx factory and stores it in the context
func (tm *PgTransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tm.trf.Begin()
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey, tx)

	err = fn(ctx)

	if tm.options.AlwaysRollback || err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
