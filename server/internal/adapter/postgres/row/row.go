package row

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/dto"
	"time"

	"github.com/google/uuid"
)

type NullAttemptItem struct {
	ID        sql.Null[uuid.UUID] `db:"id"`
	Score     sql.Null[int]       `db:"score"`
	StartedAt sql.NullTime        `db:"started_at"`
	EndedAt   sql.NullTime        `db:"ended_at"`
}

func (n *NullAttemptItem) DTO() *dto.AttemptItem {
	if !n.StartedAt.Valid {
		return nil
	}

	var endedTime *time.Time
	if n.EndedAt.Valid {
		endedTime = &n.EndedAt.Time
	}

	return &dto.AttemptItem{
		ID:        n.ID.V,
		Score:     n.Score.V,
		StartedAt: n.StartedAt.Time,
		EndedAt:   endedTime,
	}
}

func (n *NullAttemptItem) Clear() {
	*n = NullAttemptItem{}
}

type UserLastAttempt struct {
	User    sqlc.AccountProfile `db:"user"`
	Attempt NullAttemptItem     `db:"attempt"`
}

func (u *UserLastAttempt) DTO() dto.UserLastAttempt {
	return dto.UserLastAttempt{
		User: dto.User{
			ID:        u.User.AccountID,
			FirstName: u.User.FirstName,
			LastName:  u.User.LastName,
			Role:      string(u.User.Role),
		},
		LastAttempt: u.Attempt.DTO(),
	}
}

func (u *UserLastAttempt) Clear() {
	*u = UserLastAttempt{}
}

type QuizItem struct {
	Quiz    sqlc.QuizInfo       `db:"info"`
	Subject sqlc.SchoolSubject  `db:"subject"`
	Owner   sqlc.AccountProfile `db:"owner"`
}

func (q *QuizItem) Clear() {
	*q = QuizItem{}
}

type User struct {
	sqlc.AccountProfile
}

func (u *User) DTO() dto.User {
	return dto.User{
		ID:        u.AccountID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      string(u.Role),
	}
}

func (u *User) Clear() {
	*u = User{}
}

type Question struct {
	domain domain.IQuestion
	err    error
}

func (q *Question) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID              uuid.UUID             `json:"id"`
		Title           string                `json:"title"`
		Score           int                   `json:"score"`
		Type            sqlc.QuizQuestionType `json:"type"`
		SingleCorrect   string                `json:"single_correct"`
		SingleOptions   []string              `json:"single_options"`
		MultipleCorrect []string              `json:"multiple_correct"`
		MultipleOptions []string              `json:"multiple_options"`
		NumericCorrect  float32               `json:"numeric_correct"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch aux.Type {
	case sqlc.QuizQuestionTypeSingle:
		d, err := domain.RestoreSingleQuestion(
			aux.ID,
			aux.Title,
			aux.SingleCorrect,
			aux.SingleOptions,
			aux.Score,
		)
		q.domain = d
		q.err = err
	case sqlc.QuizQuestionTypeMultiple:
		d, err := domain.RestoreMultipleQuestion(
			aux.ID,
			aux.Title,
			aux.MultipleCorrect,
			aux.MultipleOptions,
			aux.Score,
		)
		q.domain = d
		q.err = err
	case sqlc.QuizQuestionTypeNumeric:
		d, err := domain.RestoreNumericQuestion(
			aux.ID,
			aux.Title,
			aux.NumericCorrect,
			aux.Score,
		)
		q.domain = d
		q.err = err
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}
