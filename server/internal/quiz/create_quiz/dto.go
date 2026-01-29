package create_quiz

import (
	"time"

	"github.com/google/uuid"
)

type Input struct {
	Title       string         `json:"title"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Summary     string         `json:"summary"`
	SubjectID   uuid.UUID      `json:"subject_id"`
	Deadline    *time.Time     `json:"deadline,omitempty"`
	MaxAttempts int            `json:"max_attempts"`
	Content     []InputContent `json:"content"`
}

type InputContent struct {
	Type    string       `json:"type"`
	Text    string       `json:"text"`
	Payload InputPayload `json:"payload"`
	Score   int          `json:"score"`
}

type InputPayload struct {
	Single   *InputSingle   `json:"single,omitempty"`
	Multiple *InputMultiple `json:"multiple,omitempty"`
	Numeric  *InputNumeric  `json:"numeric,omitempty"`
}

type InputSingle struct {
	Options []string `json:"options"`
	Correct string   `json:"correct"`
}

type InputMultiple struct {
	Options []string `json:"options"`
	Correct []string `json:"correct"`
}

type InputNumeric struct {
	Correct float64 `json:"correct"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
