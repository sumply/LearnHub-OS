package postgres

import (
	"context"
	"log/slog"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/pkg/repository/filter"
	"server/pkg/logger"

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
	log := logger.FromCtx(ctx)
	log = log.With(slog.String("id", id.String()))

	row, err := t.sqlc.GetDomainTeacher(ctx, id)
	if err != nil {
		log.WarnContext(ctx, "Failed getting row", slog.String("error", err.Error()))
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
	log := logger.FromCtx(ctx)
	log = log.With(slog.Group("teacher"))

	return t.withQueries(ctx, func(q *sqlc.Queries) error {
		err := saveDomainUser(ctx, t.Postgres, teacher.User)
		if err != nil {
			return err
		}

		for _, id := range teacher.Groups {
			err = q.InsertAccountTeacher(ctx, sqlc.InsertAccountTeacherParams{
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

func (u *User) Save(ctx context.Context, user domain.User) error {
	return u.withQueries(ctx, func(q *sqlc.Queries) error {
		err := saveDomainUser(ctx, u.Postgres, user)
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *User) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := u.sqlc.GetDomainUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        row.AccountID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      domain.UserRole(row.Role),
		CreatedAt: row.CreatedAt,
	}, nil
}

func (u *User) List(ctx context.Context, ids []uuid.UUID) ([]domain.User, error) {
	rows, err := u.sqlc.ListDomainUser(ctx, ids)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, len(rows))
	for i := range rows {
		users[i] = domain.User{
			ID:        rows[i].AccountID,
			FirstName: rows[i].FirstName,
			LastName:  rows[i].LastName,
			Role:      domain.UserRole(rows[i].Role),
			CreatedAt: rows[i].CreatedAt,
		}
	}
	return users, nil
}

func (u *User) Remove(ctx context.Context, id uuid.UUID) error {
	err := u.sqlc.DeleteAccountCredential(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Update(ctx context.Context, user domain.User) error {
	err := u.sqlc.UpdateAccountProfile(ctx, sqlc.UpdateAccountProfileParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      sqlc.AccountUserRole(user.Role),
		AccountID: user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetByCredential(ctx context.Context, filter filter.Credential) (domain.User, error) {
	row, err := u.sqlc.GetDomainUserByCredential(ctx, sqlc.GetDomainUserByCredentialParams{
		Email:   filter.Email,
		PwdHash: filter.PwdHash,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:        row.AccountID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      domain.UserRole(row.Role),
		CreatedAt: row.CreatedAt,
	}, nil
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
		err = q.InsertAccountStudent(ctx, sqlc.InsertAccountStudentParams{
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
