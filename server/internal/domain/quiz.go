package domain

import (
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID         uuid.UUID        `db:"id"`
	Title      string           `db:"title"`
	Summary    string           `db:"summary"`
	OwnerID    uuid.UUID        `db:"owner_id"`
	SubjectID  uuid.UUID        `db:"subject_id"`
	Content    []map[string]any `db:"content"`
	TotalScore int              `db:"total_score"`
	CreatedAt  time.Time        `db:"created_at"`
}
type Attempt struct {
	ID        uuid.UUID        `db:"id"`
	QuizID    uuid.UUID        `db:"quiz_id"`
	UserID    uuid.UUID        `db:"user_id"`
	Score     int              `db:"score"`
	Content   []map[string]any `db:"content"`
	StartedAt time.Time        `db:"started_at"`
	EndedAt   *time.Time       `db:"ended_at"`
}
