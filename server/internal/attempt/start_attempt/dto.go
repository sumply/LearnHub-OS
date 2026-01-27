package start_attempt

import (
	"server/internal/query"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	UserID uuid.UUID `json:"user_id"`
	QuizID uuid.UUID `json:"quiz_id"`
}

type Output struct {
	Attempt OutputAttempt
	Quiz    query.Quiz `json:"quiz"`
}

type OutputAttempt struct {
	ID        uuid.UUID `json:"id"`
	StartedAt time.Time `json:"started_at"`
}
