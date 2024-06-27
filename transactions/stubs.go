package transactions

import "context"

type TrmStub struct{}

func NewTrmStub() *TrmStub {
	return &TrmStub{}
}

func (t *TrmStub) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
