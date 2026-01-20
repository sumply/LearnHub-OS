package postgres

import (
	"context"
	"server/internal/domain"
)

func (p *Postgres) CreateUser(ctx context.Context, user *domain.User) error {
	const query = `
	INSERT INTO users 
	VALUES(
		:id,
		:first_name,
		:last_name,
		:email,
		:pwd_hash,
		:role,
		:created_at
	)`
	_, err := p.conn.NamedExecContext(ctx, query, user)

	return err
}

func (p *Postgres) CreateGroup(ctx context.Context, group *domain.Group) error {
	const query = `
	INSERT INTO groups
	VALUES(:id, :name)
	`

	_, err := p.conn.NamedExecContext(ctx, query, group)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateSubject(ctx context.Context, subject *domain.Subject) error {
	const query = `
	INSERT INTO subjects(id, name)
	VALUES(:id, :name)
	`

	_, err := p.conn.NamedExecContext(ctx, query, subject)
	if err != nil {
		return err
	}

	return nil
}
