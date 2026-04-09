package finish_attempt

import (
	"encoding/json"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Answers []RequestAnswer `json:"answers"`
}

type RequestAnswer struct {
	QuestionID uuid.UUID       `json:"question_id"`
	Answer     json.RawMessage `json:"answer"`
	a          domain.IAnswer
}

func (r *RequestAnswer) WriteSelectedAnswer(a domain.IAnswer) error {
	return a.Accept(r)
}

func (r *RequestAnswer) VisitSingleAnswer(a *domain.SingleAnswer) error {
	var selected string
	if err := json.Unmarshal(r.Answer, &selected); err != nil {
		return err
	}
	a.Selected = selected
	return nil
}

func (r *RequestAnswer) VisitMultipleAnswer(a *domain.MultipleAnswer) error {
	var selected []string
	if err := json.Unmarshal(r.Answer, &selected); err != nil {
		return err
	}
	a.Selected = selected
	return nil
}

func (r *RequestAnswer) VisitNumericAnswer(a *domain.NumericAnswer) error {
	var selected float32
	if err := json.Unmarshal(r.Answer, &selected); err != nil {
		return err
	}
	a.Selected = selected
	return nil
}

type Response struct {
	Attempt ResponseAttempt `json:"attempt"`
}

type ResponseAttempt struct {
	ID        uuid.UUID  `json:"id"`
	QuizID    uuid.UUID  `json:"quiz_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Score     int        `json:"total_score"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at"`
}
