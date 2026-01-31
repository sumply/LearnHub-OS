package finish_attempt

import (
	"server/internal/dto"

	"github.com/google/uuid"
)

type Input struct {
	Answers []InputAnswer `json:"answers"`
}

type InputAnswer struct {
	QuestionID uuid.UUID `json:"question_id"`
	Answer     any       `json:"answer"`
}

type Output struct {
	dto.FinishedAttempt
}
