package get_quiz

import "server/internal/dto"

type Output struct {
	Quizzes []dto.QuizItem `json:"quizzes"`
}
