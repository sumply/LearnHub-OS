package start_attempt

import (
	"server/internal/dto"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	UserID uuid.UUID `json:"user_id"`
}

type Output struct {
	Attempt OutputAttempt `json:"attempt"`
	Quiz    dto.Quiz      `json:"quiz"`
}

type OutputAttempt struct {
	ID        uuid.UUID `json:"id"`
	StartedAt time.Time `json:"started_at"`
}
