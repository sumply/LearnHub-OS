package create_quiz

import (
	"net/http"
	"server/internal/dto"
	"server/internal/pkg/decoder"
	"server/internal/pkg/validator"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	Title       string         `json:"title" validate:"required,quiz-title"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Summary     string         `json:"summary" validate:"required,quiz-summary"`
	SubjectID   uuid.UUID      `json:"subject_id"`
	GroupIDs    uuid.UUIDs     `json:"group_ids"`
	Deadline    *time.Time     `json:"deadline"`
	MaxAttempts int            `json:"max_attempts" validate:"required,quiz-max-attempts"`
	Questions   []dto.Question `json:"questions"`
}

func InputFromRequest(r *http.Request) (Input, error) {
	input, err := decoder.JSON[Input](r.Body)
	if err != nil {
		return Input{}, err
	}

	if err := validator.V(r.Context(), input); err != nil {
		return Input{}, err
	}

	return input, nil
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
