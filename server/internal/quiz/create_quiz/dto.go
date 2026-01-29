package create_quiz

import (
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Input struct {
	Title       string          `json:"title"`
	OwnerID     uuid.UUID       `json:"owner_id"`
	Summary     string          `json:"summary"`
	SubjectID   uuid.UUID       `json:"subject_id"`
	Deadline    *time.Time      `json:"deadline,omitempty"`
	MaxAttempts int             `json:"max_attempts"`
	Questions   []InputQuestion `json:"questions"`
}

type InputQuestion struct {
	Text    string       `json:"text"`
	Score   int          `json:"score"`
	Type    string       `json:"type"`
	Details InputDetails `json:"-"`
}

type InputDetails struct {
	Domain domain.QuestionDetails
}

func (i *InputQuestion) UnmarshalJSON(data []byte) error {
	type Alias InputQuestion
	aux := struct {
		*Alias
		RawDetails json.RawMessage `json:"details"`
	}{
		Alias: (*Alias)(i),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case "single":
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
		i.Details.Domain = &details
	case "multiple":
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
		i.Details.Domain = &details
	case "numeric":
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

		i.Details.Domain = &details
	default:
		return fmt.Errorf("type is not support")
	}

	return nil
}

type Output struct {
	ID uuid.UUID `json:"id"`
}
