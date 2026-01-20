package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (p *Postgres) UpdateProfile(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	const query = `
	UPDATE users
	SET first_name = :first_name, last_name = :last_name
	WHERE id = :id
	`
	arg := struct {
		ID        uuid.UUID `db:"id"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
	}{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}

	ext := p.selectExecuter(ctx)

	_, err := ext.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateGroup(ctx context.Context, id uuid.UUID, name string) error {
	const query = `
	UPDATE groups
	SET name = :name
	WHERE id = :id
	`
	arg := struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}{
		ID:   id,
		Name: name,
	}

	ext := p.selectExecuter(ctx)

	_, err := ext.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateSubject(ctx context.Context, id uuid.UUID, name string) error {
	const query = `
	UPDATE subjects
	SET name = :name
	WHERE id = :id
	`
	arg := struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}{
		ID:   id,
		Name: name,
	}

	ext := p.selectExecuter(ctx)

	_, err := ext.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}
