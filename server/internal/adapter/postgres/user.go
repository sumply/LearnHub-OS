package postgres

import (
	"context"
	"server/internal/adapter/postgres/sqlc"
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

func (t *Teacher) Save(ctx context.Context, teacher domain.Teacher) error {
	return t.withQueries(ctx, func(q *sqlc.Queries) error {
		err := saveDomainUser(ctx, t.Postgres, teacher.User)
		if err != nil {
			return err
		}
		for _, id := range teacher.Groups {
			err = q.InsertTeacher(ctx, sqlc.InsertTeacherParams{
				AccountID: teacher.ID,
				GroupID:   id,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type User struct {
	*Postgres
}

func NewUser(p *Postgres) *User {
	return &User{
		Postgres: p,
	}
}

func (t *User) Save(ctx context.Context, user domain.User) error {
	return t.withQueries(ctx, func(q *sqlc.Queries) error {
		err := saveDomainUser(ctx, t.Postgres, user)
		if err != nil {
			return err
		}
		return nil
	})
}

type Student struct {
	*Postgres
}

func NewStudent(p *Postgres) *Student {
	return &Student{
		Postgres: p,
	}
}

func (t *Student) Save(ctx context.Context, student domain.Student) error {
	return t.withQueries(ctx, func(q *sqlc.Queries) error {
		err := saveDomainUser(ctx, t.Postgres, student.User)
		if err != nil {
			return err
		}
		err = q.InsertStudent(ctx, sqlc.InsertStudentParams{
			AccountID: student.ID,
			GroupID:   student.Group,
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func saveDomainUser(ctx context.Context, p *Postgres, user domain.User) error {
	return p.withQueries(ctx, func(q *sqlc.Queries) error {
		err := q.InsertAccountCredential(ctx, sqlc.InsertAccountCredentialParams{
			AccountID: user.ID,
			Email:     user.Email,
			PwdHash:   user.PwdHash,
		})
		if err != nil {
			return err
		}
		err = q.InsertAccountProfile(ctx, sqlc.InsertAccountProfileParams{
			AccountID: user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      sqlc.AccountUserRole(user.Role),
			CreatedAt: user.CreatedAt,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
