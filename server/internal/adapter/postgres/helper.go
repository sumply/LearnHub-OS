package postgres

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"

	"github.com/google/uuid"
)

type QuestionJSONHelper struct {
	Domain domain.QuestionDetails
}

func (q *QuestionJSONHelper) MarshalJSON() ([]byte, error) {
	type Alias QuestionJSONHelper
	aux := struct {
		*Alias
		Variant domain.QuestionType
	}{
		Alias:   (*Alias)(q),
		Variant: q.Domain.Variant(),
	}
	return json.Marshal(&aux)
}

func (q *QuestionJSONHelper) UnmarshalJSON(data []byte) error {
	aux := struct {
		Variant    domain.QuestionType
		RawDetails json.RawMessage `json:"details"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var details domain.QuestionDetails

	switch aux.Variant {
	case domain.TypeSingleChoice:
		details = &domain.SingleChoiceQuestion{}
	case domain.TypeMultipleChoice:
		details = &domain.MultipleChoiceQuestion{}
	case domain.TypeNumeric:
		details = &domain.NumericQuestion{}
	}

	if err := json.Unmarshal(aux.RawDetails, details); err != nil {
		return err
	}

	q.Domain = details

	return nil
}

type QuizQuestionJSONAGG struct {
	ID      uuid.UUID          `json:"id"`
	QuizID  uuid.UUID          `json:"quiz_id"`
	Title   string             `json:"title"`
	Details QuestionJSONHelper `json:"details"`
	Score   int                `json:"score"`
}

func (q *QuizQuestionJSONAGG) UnmarshalJSON(data []byte) error {
	type Alias QuizQuestionJSONAGG
	aux := struct {
		*Alias
	}{
		Alias: (*Alias)(q),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if q.Details.Domain == nil {
		return fmt.Errorf("details domain is nil")
	}
	if err := q.Details.Domain.Validate(); err != nil {
		return err
	}
	return nil
}
