package find_quiz

import (
	"server/internal/dto"

	"github.com/google/uuid"
)

type Request struct {
	UserID      uuid.UUID `json:"user_id"`
	IsCompleted *bool     `json:"is_completed,omitempty"`
}

type Response struct {
	Quizzes []dto.QuizItem `json:"quizzes"`
}
