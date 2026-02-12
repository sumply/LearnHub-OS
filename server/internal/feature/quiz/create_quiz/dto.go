package create_quiz

import (
	"server/internal/dto"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	Title       string         `json:"title"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Summary     string         `json:"summary"`
	SubjectID   uuid.UUID      `json:"subject_id"`
	GroupIDs    uuid.UUIDs     `json:"group_ids"`
	Deadline    *time.Time     `json:"deadline"`
	MaxAttempts int            `json:"max_attempts"`
	Questions   []dto.Question `json:"questions"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
