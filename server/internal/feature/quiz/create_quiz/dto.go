package create_quiz

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	OwnerID     uuid.UUID         `json:"owner_id"`
	SubjectID   uuid.UUID         `json:"subject_id"`
	GroupIDs    uuid.UUIDs        `json:"group_ids"`
	Deadline    *time.Time        `json:"deadline"`
	MaxAttempts int               `json:"max_attempts"`
	Questions   []RequestQuestion `json:"questions"`
}

type RequestQuestion struct {
	Domain domain.IQuestion
	Err    error
}

func (r *RequestQuestion) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Title   string          `json:"text"`
		Score   int             `json:"score"`
		Type    string          `json:"type"`
		Details json.RawMessage `json:"details"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch aux.Type {
	case "single":
		single := &struct {
			Correct string   `json:"correct"`
			Options []string `json:"options"`
		}{}
		if err := json.Unmarshal(aux.Details, single); err != nil {
			return err
		}

		question, err := domain.NewSingleQuestion(
			aux.Title,
			single.Correct,
			single.Options,
			aux.Score,
		)

		r.Domain = question
		r.Err = err
	case "multiple":
		multiple := &struct {
			Correct []string `json:"correct"`
			Options []string `json:"options"`
		}{}
		if err := json.Unmarshal(aux.Details, multiple); err != nil {
			return err
		}

		question, err := domain.NewMultipleQuestion(
			aux.Title,
			multiple.Correct,
			multiple.Options,
			aux.Score,
		)

		r.Domain = question
		r.Err = err
	case "numeric":
		numeric := &struct {
			Correct float32 `json:"correct"`
		}{}
		if err := json.Unmarshal(aux.Details, numeric); err != nil {
			return err
		}

		question, err := domain.NewNumericQuestion(
			aux.Title,
			numeric.Correct,
			aux.Score,
		)

		r.Domain = question
		r.Err = err
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}

func (r *RequestQuestion) Unpack() (domain.IQuestion, error) {
	return r.Domain, r.Err
}

type Response struct {
	ID uuid.UUID `json:"id"`
}
