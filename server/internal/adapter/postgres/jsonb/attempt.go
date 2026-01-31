package jsonb

import (
	"server/internal/dto"

	"github.com/google/uuid"
)

type QuizAnswerAGGs []QuizAnswerAGG

func (q *QuizAnswerAGGs) ToDTOAnswers() []dto.Answer {
	answer := make([]dto.Answer, len(*q))
	for i := range answer {
		answer[i] = (*q)[i].ToDTOAnswer()
	}
	return answer
}

type QuizAnswerAGG struct {
	ID         uuid.UUID     `json:"id"`
	AttemptID  uuid.UUID     `json:"attempt_id"`
	QuestionID uuid.UUID     `json:"question_id"`
	Details    AnswerDetails `json:"details"`
	Score      int           `json:"score"`
	IsCorrect  bool          `json:"is_correct"`
}

func (q *QuizAnswerAGG) ToDTOAnswer() dto.Answer {
	return dto.Answer{
		ID:         q.ID,
		QuestionID: q.QuestionID,
		Answer:     q.Details.Answer,
		Score:      q.Score,
		IsCorrect:  q.IsCorrect,
	}
}

type AnswerDetails struct {
	Answer any `json:"answer"`
}
