package finish_attempt

import (
	"github.com/google/uuid"
)

type Input struct {
	Answers []InputAnswer `json:"answers"`
}

type InputAnswer struct {
	AttemptID uuid.UUID `json:"attempt_id"`
	Answer    any       `json:"answer"`
}

type Output struct {
	//Attempt dto.Attempt `json:"attempt"`
}
