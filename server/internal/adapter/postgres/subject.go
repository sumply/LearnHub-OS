package postgres

import (
	"context"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Subject struct {
	*Postgres
}

func NewSubject(p *Postgres) *Subject {
	return &Subject{
		Postgres: p,
	}
}

func (s *Subject) Save(ctx context.Context, subject domain.Subject) error {
	return s.sqlc.InsertSchoolSubject(ctx, sqlc.InsertSchoolSubjectParams{
		ID:   subject.ID,
		Name: subject.Name,
	})
}

func (s *Subject) Get(ctx context.Context, id uuid.UUID) (domain.Subject, error) {
	row, err := s.sqlc.GetDomainSubject(ctx, id)
	if err != nil {
		return domain.Subject{}, err
	}
	return domain.Subject{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func (s *Subject) Remove(ctx context.Context, id uuid.UUID) error {
	return s.sqlc.DeleteSchoolSubject(ctx, id)
}

func (s *Subject) Update(ctx context.Context, subject domain.Subject) error {
	err := s.sqlc.UpdateSchoolSubject(ctx, sqlc.UpdateSchoolSubjectParams{
		ID:   subject.ID,
		Name: subject.Name,
	})
	if err != nil {
		return err
	}
	return nil
}
