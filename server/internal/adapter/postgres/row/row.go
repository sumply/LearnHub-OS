package row

import (
	"database/sql"
	"server/internal/adapter/postgres/sqlc"
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

func (q *QuizItem) DTO() dto.QuizItem {
	owner := &q.Owner
	subject := &q.Subject
	quiz := &q.Quiz

	var deadline *time.Time
	if quiz.Deadline.Valid {
		deadline = &q.Quiz.Deadline.Time
	}

	return dto.QuizItem{
		ID:          quiz.QuizID,
		Title:       quiz.Title,
		Summary:     quiz.Summary,
		TotalScore:  quiz.TotalScore,
		Deadline:    deadline,
		MaxAttempts: int(quiz.MaxAttempts),
		CreatedAt:   quiz.CreatedAt,

		Owner: dto.User{
			ID:        owner.AccountID,
			FirstName: owner.FirstName,
			LastName:  owner.LastName,
			Role:      owner.AccountID.String(),
		},

		Subject: dto.Subject{
			ID:   subject.ID,
			Name: subject.Name,
		},
	}
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

type QuizLastAttempt struct {
	QuizItem
	Attempt NullAttemptItem `db:"attempt"`
}

func (q *QuizLastAttempt) DTO() dto.QuizLastAttempt {
	return dto.QuizLastAttempt{
		Quiz:        q.QuizItem.DTO(),
		LastAttempt: q.Attempt.DTO(),
	}
}

func (q *QuizLastAttempt) Clear() {
	*q = QuizLastAttempt{}
}
