package row

import (
	"database/sql"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/dto"
	"time"

	"github.com/google/uuid"
)

type UserLastAttempt struct {
	User    sqlc.AccountProfile `db:"user"`
	Attempt struct {
		ID        sql.Null[uuid.UUID] `db:"id"`
		Score     sql.Null[int]       `db:"score"`
		StartedAt sql.NullTime        `db:"started_at"`
		EndedAt   sql.NullTime        `db:"ended_at"`
	}
}

func (u *UserLastAttempt) DTO() dto.UserLastAttempt {
	var attemptItem *dto.AttemptItem
	if u.Attempt.StartedAt.Valid {
		var endedAt *time.Time
		if u.Attempt.EndedAt.Valid {
			endedAt = &u.Attempt.EndedAt.Time
		}
		attemptItem = &dto.AttemptItem{
			ID:        u.Attempt.ID.V,
			StartedAt: u.Attempt.StartedAt.Time,
			EndedAt:   endedAt,
			Score:     u.Attempt.Score.V,
		}
	}
	return dto.UserLastAttempt{
		User: dto.User{
			ID:        u.User.AccountID,
			FirstName: u.User.FirstName,
			LastName:  u.User.LastName,
			Role:      string(u.User.Role),
		},
		LastAttempt: attemptItem,
	}
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
