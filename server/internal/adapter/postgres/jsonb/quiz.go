package jsonb

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"

	"github.com/google/uuid"
)

type QuestionDetails struct {
	Domain domain.QuestionDetails
}

func (q *QuestionDetails) ToDTOQuestionDetails() dto.QuestionDetails {
	return dto.QuestionDetails{
		Domain: q.Domain,
	}
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

type QuizQuestionAGGs []QuizQuestionAGG

func (q *QuizQuestionAGGs) ToDTOQuestions() []dto.Question {
	questions := make([]dto.Question, len(*q))
	for i := range questions {
		questions[i] = (*q)[i].ToDTOQuestion()
	}
	return questions
}

func (q *QuizQuestionAGGs) ToDomainQuestions() []domain.Question {
	questions := make([]domain.Question, len(*q))
	for i := range questions {
		questions[i] = (*q)[i].ToDomainQuestion()
	}
	return questions
}

type QuizQuestionAGG struct {
	ID      uuid.UUID       `json:"id"`
	QuizID  uuid.UUID       `json:"quiz_id"`
	Title   string          `json:"title"`
	Details QuestionDetails `json:"details"`
	Score   int             `json:"score"`
}

func (q *QuizQuestionAGG) ToDTOQuestion() dto.Question {
	return dto.Question{
		ID:      q.ID,
		Text:    q.Title,
		Score:   q.Score,
		Details: q.Details.ToDTOQuestionDetails(),
	}
}

func (q *QuizQuestionAGG) ToDomainQuestion() domain.Question {
	return domain.Question{
		ID:      q.ID,
		QuizID:  q.QuizID,
		Text:    q.Title,
		Score:   q.Score,
		Details: q.Details.Domain,
	}
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
