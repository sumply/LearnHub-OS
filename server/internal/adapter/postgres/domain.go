package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	*Postgres
}

func (d *Domain) AddGroup(ctx context.Context, group *domain.Group) error {
	err := d.sqlc.InsertSchoolGroup(ctx, sqlc.InsertSchoolGroupParams{
		ID:   group.ID,
		Name: group.Name,
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *Domain) GetUser(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user := d.tables.AccountProfile
	ds := d.goqu.From(user).
		Select(user.All()).
		Where(
			user.Col("account_id").Eq(userID),
		)

	var row sqlc.AccountProfile
	ok, err := ds.Executor().ScanStructContext(ctx, &row)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, repository.NewNotFoundError()
	}

	return domain.User{
		ID:        row.AccountID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Role:      domain.UserRole(row.Role),
		CreatedAt: row.CreatedAt,
	}, nil
}

func (d *Domain) ListUserByIDs(ctx context.Context, userIDs uuid.UUIDs) ([]domain.User, error) {
	user := d.tables.AccountProfile
	ds := d.goqu.From(user).
		Select(user.All()).
		Where(
			user.Col("account_id").In(userIDs),
		)

	var rows []sqlc.AccountProfile
	if err := ds.Executor().ScanStructsContext(ctx, &rows); err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, domain.User{
			ID:        row.AccountID,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Role:      domain.UserRole(row.Role),
			CreatedAt: row.CreatedAt,
		})
	}

	return users, nil
}

func (d *Domain) AddTeacherToGroup(ctx context.Context, groupID uuid.UUID, teacherID uuid.UUID) error {
	teacher := d.tables.AccountTeacher

	row := sqlc.AccountTeacher{
		GroupID:   groupID,
		AccountID: teacherID,
	}

	ds := d.goqu.Insert(teacher).Rows(&row)
	if _, err := ds.Executor().ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (d *Domain) AddStudentsToGroup(ctx context.Context, groupID uuid.UUID, studentIDs uuid.UUIDs) error {
	student := d.tables.AccountStudent

	rows := make([]sqlc.AccountStudent, 0, len(studentIDs))
	for _, id := range studentIDs {
		rows = append(rows, sqlc.AccountStudent{
			GroupID:   groupID,
			AccountID: id,
		})
	}

	ds := d.goqu.Insert(student).Rows(rows)
	if _, err := ds.Executor().ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DomainQuiz(ctx context.Context, id uuid.UUID) (domain.Quiz, error) {
	row, err := p.sqlc.GetDomainQuiz(ctx, id)
	if err != nil {
		return domain.Quiz{}, fmt.Errorf("getting domain quiz: %w", err)
	}

	var deadline *time.Time
	if row.QuizDeadline.Valid {
		deadline = &row.QuizDeadline.Time
	}

	var questions jsonb.QuizQuestionAGGs
	if err := json.Unmarshal(row.QuizQuestions, &questions); err != nil {
		return domain.Quiz{}, err
	}

	quiz := domain.Quiz{
		ID:          row.QuizID,
		Title:       row.QuizTitle,
		Summary:     row.QuizSummary,
		OwnerID:     row.QuizOwnerID,
		SubjectID:   row.QuizSubjectID,
		Deadline:    deadline,
		Questions:   questions.ToDomainQuestions(),
		MaxAttempts: int(row.QuizMaxAttempts),
		TotalScore:  int(row.QuizTotalScore),
		CreatedAt:   row.QuizCreatedAt,
	}

	return quiz, nil
}

func (p *Postgres) DomainAttempt(ctx context.Context, id uuid.UUID) (domain.Attempt, error) {
	row, err := p.sqlc.GetDomainAttempt(ctx)
	if err != nil {
		return domain.Attempt{}, fmt.Errorf("getting domain attempt: %w", err)
	}

	var answers jsonb.QuizAnswerAGGs
	if err := json.Unmarshal(row.AttemptAnswers, &answers); err != nil {
		return domain.Attempt{}, err
	}

	var endedAt *time.Time
	if row.AttemptEndendAt.Valid {
		endedAt = &row.AttemptEndendAt.Time
	}

	return domain.Attempt{
		ID:        row.AttemptID,
		QuizID:    row.AttemptQuizID,
		UserID:    row.AttemptUserID,
		Answers:   answers.ToDomainAnswers(),
		Score:     int(row.AttemptScoreID),
		StartedAt: row.AttemptStartedAt,
		EndedAt:   endedAt,
	}, nil
}
