package dto

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type QuizItem struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Owner       User       `json:"owner"`
	Subject     Subject    `json:"subject"`
	TotalScore  int        `json:"total_score"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	MaxAttempts int        `json:"max_attempts"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Quiz struct {
	QuizItem
	Content []Question `json:"content"`
}

type Question struct {
	ID      uuid.UUID       `json:"id"`
	Text    string          `json:"text"`
	Score   int             `json:"score"`
	Details QuestionDetails `json:"-"`
}

type QuestionDetails struct {
	Domain domain.QuestionDetails
}

func (q *Question) UnmarshalJSON(data []byte) error {
	type Alias Question
	aux := struct {
		*Alias
		Type       domain.QuestionType `json:"type"`
		RawDetails json.RawMessage     `json:"details"`
	}{
		Alias: (*Alias)(q),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case domain.TypeSingleChoice:
		single := struct {
			Options []string `json:"options"`
			Correct string   `json:"correct"`
		}{}
		if err := json.Unmarshal(aux.RawDetails, &single); err != nil {
			return err
		}
		details, err := domain.NewSingleChoiceQuestion(single.Options, single.Correct)
		if err != nil {
			return err
		}
		q.Details.Domain = &details
	case domain.TypeMultipleChoice:
		multiple := struct {
			Options []string `json:"options"`
			Correct []string `json:"correct"`
		}{}
		if err := json.Unmarshal(aux.RawDetails, &multiple); err != nil {
			return err
		}

		details, err := domain.NewMultipleChoiceQuestion(multiple.Options, multiple.Correct)
		if err != nil {
			return err
		}
		q.Details.Domain = &details
	case domain.TypeNumeric:
		numeric := struct {
			Correct float64 `json:"correct"`
		}{}
		if err := json.Unmarshal(aux.RawDetails, &numeric); err != nil {
			return err
		}

		details, err := domain.NewNumericQuestion(numeric.Correct)
		if err != nil {
			return err
		}

		q.Details.Domain = &details
	default:
		return fmt.Errorf("type is not support")
	}

	return nil
}

func (q *Question) MarshalJSON() ([]byte, error) {
	type Alias Question
	aux := struct {
		*Alias
		Type    domain.QuestionType `json:"type"`
		Details any                 `json:"details"`
	}{
		Alias: (*Alias)(q),
	}

	if q.Details.Domain == nil {
		return nil, fmt.Errorf("domain is nil")
	}

	switch q.Details.Domain.Variant() {
	case domain.TypeSingleChoice:
		details, ok := q.Details.Domain.(*domain.SingleChoiceQuestion)
		if !ok {
			return nil, fmt.Errorf("invalid details type")
		}
		aux.Details = &struct {
			Options []string `json:"options"`
			Correct string   `json:"correct"`
		}{
			Options: details.Options,
			Correct: details.Correct,
		}
	case domain.TypeMultipleChoice:
		details, ok := q.Details.Domain.(*domain.MultipleChoiceQuestion)
		if !ok {
			return nil, fmt.Errorf("invalid details type")
		}
		aux.Details = &struct {
			Options []string `json:"options"`
			Correct []string `json:"correct"`
		}{
			Options: details.Options,
			Correct: details.Correct,
		}
	case domain.TypeNumeric:
		details, ok := q.Details.Domain.(*domain.NumericQuestion)
		if !ok {
			return nil, fmt.Errorf("invalid details type")
		}
		aux.Details = &struct {
			Correct float64 `json:"correct"`
		}{
			Correct: details.Correct,
		}
	default:
		return nil, fmt.Errorf("incorrected QuestionType")
	}

	aux.Type = q.Details.Domain.Variant()

	return json.Marshal(aux)
}
