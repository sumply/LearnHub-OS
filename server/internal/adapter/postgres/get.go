package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/dto"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

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
				User:     user,
				Groups:   rows[i].TGroupIds,
				Subjects: rows[i].TSubjectIds,
			}
			users[i] = &teacher
		default:
			users[i] = user
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

	var deadline *time.Time
	if row.QuizDeadline.Valid {
		deadline = &row.QuizDeadline.Time
	}

	var helper []jsonb.QuizQuestionAGG
	if err := json.Unmarshal(row.QuizQuestions, &helper); err != nil {
		return dto.Quiz{}, err
	}

	questions := make([]dto.Question, len(helper))
	for i := range questions {
		questions[i] = dto.Question{
			Text:  helper[i].Title,
			Score: helper[i].Score,
			Details: dto.QuestionDetails{
				Domain: helper[i].Details.Domain,
			},
		}
	}

	quiz := dto.Quiz{
		QuizItem: dto.QuizItem{
			ID:      row.QuizID,
			Title:   row.QuizTitle,
			Summary: row.QuizSummary,
			Owner: dto.User{
				ID:        row.OwnerID,
				FirstName: row.OwnerFirstName,
				LastName:  row.OwnerLastName,
				Role:      string(row.OwnerRole),
			},
			Subject: dto.Subject{
				ID:   row.SubjectID,
				Name: row.SubjectName,
			},
			TotalScore:  int(row.QuizTotalScore),
			Deadline:    deadline,
			MaxAttempts: int(row.QuizMaxAttempts),
			CreatedAt:   row.QuizCreatedAt,
		},
		Content: questions,
	}

	return quiz, nil
}

func (p *Postgres) QuizWithoutAnswers(ctx context.Context, id uuid.UUID) (dto.Quiz, error) {
	row, err := p.sqlc.GetQuiz(ctx, id)
	if err != nil {
		return dto.Quiz{}, err
	}

	var deadline *time.Time
	if row.QuizDeadline.Valid {
		deadline = &row.QuizDeadline.Time
	}

	fmt.Println(string(row.QuizQuestions))

	var content []dto.Question
	json.Unmarshal(row.QuizQuestions, &content)

	quiz := dto.Quiz{
		QuizItem: dto.QuizItem{
			ID:      row.QuizID,
			Title:   row.QuizTitle,
			Summary: row.QuizSummary,
			Owner: dto.User{
				ID:        row.OwnerID,
				FirstName: row.OwnerFirstName,
				LastName:  row.OwnerLastName,
				Role:      string(row.OwnerRole),
			},
			Subject: dto.Subject{
				ID:   row.SubjectID,
				Name: row.SubjectName,
			},
			TotalScore:  int(row.QuizTotalScore),
			Deadline:    deadline,
			MaxAttempts: int(row.QuizMaxAttempts),
			CreatedAt:   row.QuizCreatedAt,
		},
		Content: content,
	}

	return quiz, nil
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

func (p *Postgres) UpdateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	return nil
}
