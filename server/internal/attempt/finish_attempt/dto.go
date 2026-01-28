package finish_attempt

import "github.com/google/uuid"

type Input struct {
	Answers []InputAnswer `json:"answers"`
}

type InputAnswer struct {
	AnswerID uuid.UUID `json:"answer_id"`
	Answer   any       `json:"answer"`
}

type Output struct{}
