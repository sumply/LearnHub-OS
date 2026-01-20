package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (p *Postgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM users
	WHERE id=$1
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM groups
	WHERE id = $1
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM subjects
	WHERE id = $1
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
