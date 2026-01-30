package jsonb

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"

	"github.com/google/uuid"
)

type QuestionDetails struct {
	Domain domain.QuestionDetails
}

func (q QuestionDetails) MarshalJSON() ([]byte, error) {
	type Alias QuestionDetails
	aux := struct {
		Details domain.QuestionDetails `json:"details"`
		Variant domain.QuestionType    `json:"variant"`
	}{
		Details: q.Domain,
		Variant: q.Domain.Variant(),
	}
	return json.Marshal(&aux)
}

func (q *QuestionDetails) UnmarshalJSON(data []byte) error {
	aux := struct {
		Variant    domain.QuestionType `json:"variant"`
		RawDetails json.RawMessage     `json:"details"`
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

type QuizQuestionAGG struct {
	ID      uuid.UUID       `json:"id"`
	QuizID  uuid.UUID       `json:"quiz_id"`
	Title   string          `json:"title"`
	Details QuestionDetails `json:"details"`
	Score   int             `json:"score"`
}

func (q *QuizQuestionAGG) UnmarshalJSON(data []byte) error {
	type Alias QuizQuestionAGG
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
