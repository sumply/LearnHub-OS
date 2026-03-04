package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (p *Postgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user := p.tables.AccountCredential

	ds := p.goqu.From(user).
		Delete().
		Where(user.Col("account_id").Eq(id))

	_, err := ds.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM school.group
	WHERE id = $1
	`

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM school.subject
	WHERE id = $1
	`

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM quiz.info
	WHERE id = $1
	`

	_, err := p.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
