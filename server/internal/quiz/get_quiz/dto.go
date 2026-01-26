package get_quiz

import "server/internal/query"

type Output struct {
	Quizzes []query.QuizItem `json:"quizzes"`
}
