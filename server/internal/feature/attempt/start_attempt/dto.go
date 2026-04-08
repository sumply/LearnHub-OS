package start_attempt

import "github.com/google/uuid"

type Request struct {
	QuizID uuid.UUID `json:"quiz_id"`
	UserID uuid.UUID `json:"user_id"`
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
