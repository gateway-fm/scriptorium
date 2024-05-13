package transactions

import (
	"context"

	"github.com/go-pg/pg/v10/orm"
)

type contextKey string

const txKey contextKey = "tx"

type Transaction interface {
	orm.DB
}

type TransactionFactory interface {
	Transaction(ctx context.Context) Transaction
}

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
