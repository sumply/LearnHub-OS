package get_quizzes

import "server/internal/dto"

type Output struct {
	Quizzes []dto.FinishedAttempt `json:"quizzes"`
}
