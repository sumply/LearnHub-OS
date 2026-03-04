package postgres

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Teacher struct {
	*Postgres
}

func NewTeacher(p *Postgres) *Teacher {
	return &Teacher{
		Postgres: p,
	}
}

func (t *Teacher) Get(ctx context.Context, id uuid.UUID) (domain.Teacher, error) {
	row, err := t.sqlc.GetDomainTeacher(ctx, id)
	if err != nil {
		return domain.Teacher{}, err
	}

	return domain.Teacher{
		User: domain.User{
			ID:        row.AccountID,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Role:      domain.UserRole(row.Role),
			CreatedAt: row.CreatedAt,
		},
		Groups: row.GroupIds,
	}, nil
}
