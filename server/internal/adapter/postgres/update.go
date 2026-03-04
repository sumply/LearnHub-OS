package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (p *Postgres) UpdateProfile(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	user := p.tables.AccountProfile

	ds := p.goqu.From(user).
		Update().
		Set(goqu.Record{
			"first_name": firstName,
			"last_name":  lastName,
		}).
		Where(
			user.Col("account_id").Eq(id),
		)

	_, err := ds.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateGroup(ctx context.Context, id uuid.UUID, name string) error {
	const query = `
	UPDATE school.group
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

	ext, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = ext.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateSubject(ctx context.Context, id uuid.UUID, name string) error {
	const query = `
	UPDATE school.subject
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

	ext, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = ext.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	var endedAt sql.NullTime
	if attempt.EndedAt != nil {
		endedAt.Time = *attempt.EndedAt
		endedAt.Valid = true
	}
	err := p.sqlc.UpdateQuizAttempt(ctx, sqlc.UpdateQuizAttemptParams{
		Score:   attempt.Score,
		EndedAt: endedAt,
		ID:      attempt.ID,
	})
	if err != nil {
		return err
	}

	for _, answer := range attempt.Answers {
		details, err := json.Marshal(jsonb.AnswerDetails{Answer: answer.Answer})
		if err != nil {
			return err
		}
		err = p.sqlc.UpdateQuizAnswer(ctx, sqlc.UpdateQuizAnswerParams{
			Details:   details,
			Score:     answer.Score,
			IsCorrect: answer.IsCorrect,
			ID:        answer.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
