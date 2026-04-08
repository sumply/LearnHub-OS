package list_attempts

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Attempts []ResponseAttempt `json:"attempt"`
}

type ResponseAttempt struct {
	ID        uuid.UUID  `json:"id"`
	QuizID    uuid.UUID  `json:"quiz_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Score     int        `json:"score"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}
