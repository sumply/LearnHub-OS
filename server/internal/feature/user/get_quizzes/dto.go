package get_quizzes

import "server/internal/dto"

type Output struct {
	Quizzes []dto.QuizLastAttempt `json:"quizzes"`
}
