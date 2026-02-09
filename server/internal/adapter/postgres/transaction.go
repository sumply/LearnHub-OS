package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	conn *sqlx.DB
}

func (t *Transaction) Execute(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	isComplete := false

	defer func() {
		if !isComplete {
			tx.Rollback()
		}
	}()

	ctx = ctxWithTx(ctx, tx)

	err = fn(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	isComplete = true

	return nil
}

type ctxKey string

const txKey ctxKey = "tx"

func txFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sqlx.Tx)
	return tx, ok
}

func ctxWithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
