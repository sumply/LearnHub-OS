package postgres

import (
	"context"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"
	"server/internal/pkg/repository/filter"
	"time"

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

func (p *Postgres) Quiz(ctx context.Context, id uuid.UUID) (dto.Quiz, error) {
	row, err := p.sqlc.GetQuiz(ctx, id)
	if err != nil {
		return dto.Quiz{}, err
	}

	return quizRow{&row}.ToDTOQuiz()
}

type quizRow struct {
	*sqlc.GetQuizRow
}

func (g quizRow) ToDTOQuiz() (dto.Quiz, error) {
	var deadline *time.Time
	if g.QuizDeadline.Valid {
		deadline = &g.QuizDeadline.Time
	}

	var questions jsonb.QuizQuestionAGGs
	if err := json.Unmarshal(g.QuizQuestions, &questions); err != nil {
		return dto.Quiz{}, err
	}

	return dto.Quiz{
		QuizItem: dto.QuizItem{
			ID:      g.QuizID,
			Title:   g.QuizTitle,
			Summary: g.QuizSummary,
			Owner: dto.User{
				ID:        g.OwnerID,
				FirstName: g.OwnerFirstName,
				LastName:  g.OwnerLastName,
				Role:      string(g.OwnerRole),
			},
			Subject: dto.Subject{
				ID:   g.SubjectID,
				Name: g.SubjectName,
			},
			TotalScore:  int(g.QuizTotalScore),
			Deadline:    deadline,
			MaxAttempts: int(g.QuizMaxAttempts),
			CreatedAt:   g.QuizCreatedAt,
		},
		Content: questions.ToDTOQuestions(),
	}, nil
}

func (p *Postgres) FinishedAttempt(ctx context.Context, id uuid.UUID) (dto.FinishedAttempt, error) {
	row, err := p.sqlc.GetFinishedAttempt(ctx, id)
	if err != nil {
		return dto.FinishedAttempt{}, err
	}

	return finishedAttemptRow{&row}.ToDTOFinishedAttempt()
}

type finishedAttemptRow struct {
	*sqlc.GetFinishedAttemptRow
}

func (g finishedAttemptRow) ToDTOFinishedAttempt() (dto.FinishedAttempt, error) {
	var questions jsonb.QuizQuestionAGGs
	if err := json.Unmarshal(g.QuizQuestions, &questions); err != nil {
		return dto.FinishedAttempt{}, err
	}

	var answers jsonb.QuizAnswerAGGs
	if err := json.Unmarshal(g.AttemptAnswers, &answers); err != nil {
		return dto.FinishedAttempt{}, err
	}

	var deadline *time.Time
	if g.QuizDeadline.Valid {
		deadline = &g.QuizDeadline.Time
	}
	return dto.FinishedAttempt{
		Attempt: dto.Attempt{
			ID: g.AttemptID,
			User: dto.User{
				ID:        g.UserID,
				FirstName: g.UserFirstName,
				LastName:  g.UserLastName,
				Role:      string(g.UserRole),
			},
			Answers:   answers.ToDTOAnswers(),
			Score:     int(g.AttemptScore),
			StartedAt: g.AttemptStartedAt,
			EndedAt:   g.AttemptEndedAt.Time,
		},
		Quiz: dto.Quiz{
			QuizItem: dto.QuizItem{
				ID:      g.QuizID,
				Title:   g.QuizTitle,
				Summary: g.QuizSummary,
				Owner: dto.User{
					ID:        g.OwnerID,
					FirstName: g.OwnerFirstName,
					LastName:  g.OwnerLastName,
					Role:      string(g.OwnerRole),
				},
				Subject: dto.Subject{
					ID:   g.SubjectID,
					Name: g.SubjectName,
				},
				TotalScore:  int(g.QuizTotalScore),
				Deadline:    deadline,
				MaxAttempts: int(g.QuizMaxAttempts),
				CreatedAt:   g.QuizCreatedAt,
			},
			Content: questions.ToDTOQuestions(),
		},
	}, nil
}

func (p *Postgres) QuizItems(ctx context.Context) ([]dto.QuizItem, error) {
	rows, err := p.sqlc.GetQuizItem(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.QuizItem, len(rows))
	for i := range items {
		var deadline *time.Time
		if rows[i].QuizDeadline.Valid {
			deadline = &rows[i].QuizDeadline.Time
		}
		items[i] = dto.QuizItem{
			ID:      rows[i].QuizID,
			Title:   rows[i].QuizTitle,
			Summary: rows[i].QuizSummary,
			Owner: dto.User{
				ID:        rows[i].OwnerID,
				FirstName: rows[i].OwnerFirstName,
				LastName:  rows[i].OwnerLastName,
				Role:      string(rows[i].OwnerRole),
			},
			Subject: dto.Subject{
				ID:   rows[i].SubjectID,
				Name: rows[i].SubjectName,
			},
			TotalScore:  int(rows[i].QuizTotalScore),
			Deadline:    deadline,
			MaxAttempts: int(rows[i].QuizMaxAttempts),
			CreatedAt:   rows[i].OwnerCreatedAt,
		}
	}

	return items, nil
}
