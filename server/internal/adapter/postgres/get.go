package postgres

import (
	"context"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"
	"server/internal/pkg/repository/filter"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (p *Postgres) User(ctx context.Context, userID uuid.UUID) (dto.User, error) {
	user := p.tables.AccountProfile
	ds := p.goqu.From(user).
		Select(user.All()).
		Where(
			user.Col("account_id").Eq(userID),
		)
	var row sqlc.AccountProfile
	ok, err := ds.Executor().ScanStructContext(ctx, &row)
	if err != nil {
		return dto.User{}, err
	}
	if !ok {
		return dto.User{}, repository.NewNotFoundError()
	}

	return dto.User{
		ID:        row.AccountID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      string(row.Role),
	}, nil
}

func (p *Postgres) buildUserByEmailQuery(filter filter.Credential) (*goqu.SelectDataset, *sqlc.AccountProfile) {
	credential := p.tables.AccountCredential
	profile := p.tables.AccountProfile

	profileOn := goqu.On(profile.Col("account_id").Eq(credential.Col("account_id")))

	ds := p.goqu.From(credential).
		Select(profile.All()).
		Join(profile, profileOn).
		Where(
			credential.Col("email").Eq(filter.Email),
			credential.Col("pwd_hash").Eq(filter.PwdHash),
		)

	return ds, &sqlc.AccountProfile{}
}

func (p *Postgres) DetailedUsers(ctx context.Context) ([]domain.UserAggregate, error) {
	rows, err := p.sqlc.GetDetailedUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]domain.UserAggregate, len(rows))
	for i := range rows {
		user := domain.User{
			ID:        rows[i].AccountID,
			FirstName: rows[i].FirstName,
			LastName:  rows[i].LastName,
			Role:      domain.UserRole(rows[i].Role),
			CreatedAt: rows[i].CreatedAt,
		}
		switch rows[i].Role {
		case sqlc.AccountUserRoleStudent:
			student := domain.Student{
				User:  user,
				Group: rows[i].SGroupID.UUID,
			}
			users[i] = &student
		case sqlc.AccountUserRoleTeacher:
			teacher := domain.Teacher{
				User:   user,
				Groups: rows[i].TGroupIds,
			}
			users[i] = &teacher
		default:
			users[i] = &user
		}
	}

	return users, nil
}

func (p *Postgres) Groups(ctx context.Context) ([]domain.Group, error) {
	const query = `
	SELECT 
		id,
		name
	FROM school.group
	`

	var groups []domain.Group
	err := p.conn.SelectContext(ctx, &groups, query)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (p *Postgres) Subjects(ctx context.Context) ([]domain.Subject, error) {
	const query = `
	SELECT id, name
	FROM school.subject
	`

	var subjects []domain.Subject
	err := p.conn.SelectContext(ctx, &subjects, query)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}
