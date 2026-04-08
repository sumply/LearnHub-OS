package get_quiz

import (
	"encoding/json"
	"os"
	"server/internal/domain"
	"server/internal/pkg/encoder"
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID          uuid.UUID           `json:"id"`
	Title       string              `json:"title"`
	Summary     string              `json:"summary"`
	OwnerID     uuid.UUID           `json:"owner_id"`
	SubjectID   uuid.UUID           `json:"subject_id"`
	Questions   []*ResponseQuestion `json:"questions"`
	GroupIDs    uuid.UUIDs          `json:"group_ids"`
	Deadline    *time.Time          `json:"deadline"`
	MaxAttempts int                 `json:"max_attempts"`
	TotalScore  int                 `json:"total_score"`
	CreatedAt   time.Time           `json:"created_at"`
}

type ResponseQuestion struct {
	data map[string]any
}

func NewResponseQuestion(q domain.IQuestion) (*ResponseQuestion, error) {
	rq := &ResponseQuestion{}
	if err := q.Accept(rq); err != nil {
		return nil, err
	}
	encoder.JSON(os.Stdout, rq)
	return rq, nil
}

func (r *ResponseQuestion) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.data)
}

func (r *ResponseQuestion) VisitSingleQuestion(single *domain.SingleQuestion) error {
	r.data = map[string]any{
		"id":      single.ID(),
		"title":   single.Title,
		"type":    "single",
		"score":   single.Score(),
		"options": single.Options,
		"correct": single.Correct,
	}
	return nil
}

func (r *ResponseQuestion) VisitMultipleQuestion(multiple *domain.MultipleQuestion) error {
	r.data = map[string]any{
		"id":      multiple.ID(),
		"title":   multiple.Title,
		"type":    "multiple",
		"score":   multiple.Score(),
		"options": multiple.Options,
		"correct": multiple.Correct,
	}
	return nil
}

func (r *ResponseQuestion) VisitNumericQuestion(numeric *domain.NumericQuestion) error {
	r.data = map[string]any{
		"id":      numeric.ID(),
		"title":   numeric.Title,
		"type":    "numeric",
		"score":   numeric.Score(),
		"correct": numeric.Correct,
	}
	return nil
}
