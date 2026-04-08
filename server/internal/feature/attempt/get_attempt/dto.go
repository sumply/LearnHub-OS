package get_attempt

import (
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	Attempt ResponseAttempt `json:"attempt"`
}

type ResponseAttempt struct {
	ID        uuid.UUID        `json:"id"`
	QuizID    uuid.UUID        `json:"quiz_id"`
	UserID    uuid.UUID        `json:"user_id"`
	Score     int              `json:"score"`
	Answers   []ResponseAnswer `json:"answers"`
	StartedAt time.Time        `json:"started_at"`
	EndedAt   *time.Time       `json:"ended_at"`
}

type ResponseAnswer struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Answer     any       `json:"answer"`
	Score      int       `json:"score"`
	IsCorrect  bool      `json:"is_correct"`
}

func (r *ResponseAnswer) VisitSingleAnswer(answer *domain.SingleAnswer) error {
	r.mapIAnswer(answer)
	r.Answer = answer.Selected
	return nil
}

func (r *ResponseAnswer) VisitMultipleAnswer(answer *domain.MultipleAnswer) error {
	r.mapIAnswer(answer)
	r.Answer = answer.Selected
	return nil
}

func (r *ResponseAnswer) VisitNumericAnswer(answer *domain.NumericAnswer) error {
	r.mapIAnswer(answer)
	r.Answer = answer.Selected
	return nil
}

func (r *ResponseAnswer) mapIAnswer(answer domain.IAnswer) {
	r.ID = answer.ID()
	r.QuestionID = answer.QuestionID()
	r.Score = answer.Score()
	r.IsCorrect = answer.IsCorrect()
}
