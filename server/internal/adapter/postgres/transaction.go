package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ctxKey string

const (
	txKey ctxKey = "tx"
)

type Transaction struct {
	*Postgres
}

func (t *Transaction) With(ctx context.Context, f func(context.Context) error) error {
	tx, err := t.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	ctx = context.WithValue(ctx, txKey, tx)
	if err := f(ctx); err != nil {
		return err
	}

	return tx.Commit()
}

func txFromCtx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}
