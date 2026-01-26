package create_quiz

import (
	"server/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateQuestionPayload(t *testing.T) {
	tests := []struct {
		Name         string
		Payload      InputPayload
		ExpectedType any
	}{
		{
			Name: "single",
			Payload: InputPayload{
				Single: &InputSingle{
					Options: []string{"correct", "incorrect"},
					Correct: 0,
				},
			},
			ExpectedType: &domain.SingleChoiceQuestion{},
		},
		{
			Name: "multiple",
			Payload: InputPayload{
				Multiple: &InputMultiple{
					Options: []string{"correct", "incorrect"},
					Correct: []int{0, 1},
				},
			},
			ExpectedType: &domain.MultipleChoiceQuestion{},
		},
		{
			Name: "numeric",
			Payload: InputPayload{
				Numeric: &InputNumeric{
					Correct: 0.1,
				},
			},
			ExpectedType: &domain.NumericQuestion{},
		},
		{
			Name:         "invalid",
			Payload:      InputPayload{},
			ExpectedType: nil,
		},
	}
	usecase := UseCase{}
	for i := range tests {
		t.Run(tests[i].Name, func(t *testing.T) {
			result := usecase.createQuestionPayload(&tests[i].Payload)
			if !assert.IsType(t, tests[i].ExpectedType, result) {
				t.Errorf("unexpected type; want: %T; got: %T", tests[i].ExpectedType, result)
			}
		})
	}
}
