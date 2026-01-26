package domain

import "testing"

func TestSingleChoiceQuestion(t *testing.T) {
	tests := []struct {
		Name        string
		question    SingleChoiceQuestion
		ExpectError bool
	}{
		{
			Name: "Empty",
			question: SingleChoiceQuestion{
				Options: nil,
				Correct: 0,
			},
			ExpectError: true,
		},
		{
			Name: "Overflow",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: 3,
			},
			ExpectError: true,
		},
		{
			Name: "Underflow",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: -1,
			}, ExpectError: true,
		},
		{
			Name: "Correct",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: 0,
			},
			ExpectError: false,
		},
		{
			Name: "Correct",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: 1,
			},
			ExpectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].Name, func(t *testing.T) {
			err := tests[i].question.Validate()
			if tests[i].ExpectError && err == nil {
				t.Errorf("expect error but got nil")
			}
			if !tests[i].ExpectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMultipleChoiceQuestion(t *testing.T) {
	tests := []struct {
		Name        string
		question    MultipleChoiceQuestion
		ExpectError bool
	}{
		{
			Name: "Options is empty",
			question: MultipleChoiceQuestion{
				Options: nil,
				Correct: []int{0},
			},
			ExpectError: true,
		},
		{
			Name: "Correct is empty",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1"},
				Correct: nil,
			},
			ExpectError: true,
		},
		{
			Name: "Correct contains double1",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1"},
				Correct: []int{0, 0},
			},
			ExpectError: true,
		},
		{
			Name: "Correct contains double2",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1"},
				Correct: []int{0, 1, 0},
			},
			ExpectError: true,
		},
		{
			Name: "Overflow",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []int{0, 1, 2},
			},
			ExpectError: true,
		},
		{
			Name: "Underflow",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []int{-1, 0, 1},
			},
			ExpectError: true,
		},
		{
			Name: "Correct",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []int{0, 1},
			},
			ExpectError: false,
		},
	}
	for i := range tests {
		t.Run(tests[i].Name, func(t *testing.T) {
			err := tests[i].question.Validate()
			if tests[i].ExpectError && err == nil {
				t.Errorf("expect error but got nil")
			}
			if !tests[i].ExpectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNumericQuestion(t *testing.T) {
	tests := []struct {
		name        string
		question    NumericQuestion
		expectError bool
	}{
		{
			name: "Less 0",
			question: NumericQuestion{
				Correct: -34,
			},
			expectError: true,
		},
		{
			name: "Correct in border",
			question: NumericQuestion{
				Correct: 0,
			},
			expectError: false,
		},
		{
			name: "Correct",
			question: NumericQuestion{
				Correct: 123,
			},
			expectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			err := tests[i].question.Validate()
			if tests[i].expectError && err == nil {
				t.Errorf("expect error but got nil")
			}
			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestQuestionAggregate(t *testing.T) {
	tests := []struct {
		name        string
		aggregate   QuestionAggregate
		expectError bool
	}{
		{
			name: "Invalid type",
			aggregate: QuestionAggregate{
				Type: "invalid",
			},
			expectError: true,
		},
		{
			name: "Invalid payload as numeric",
			aggregate: QuestionAggregate{
				Type: TypeNumeric,
				Payload: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: 0,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid payload as single",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Payload: &MultipleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: []int{0},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid payload as multiple",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Payload: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: 0,
				},
			},
			expectError: true,
		},
		{
			name: "Payload is nil",
			aggregate: QuestionAggregate{
				Type:    TypeNumeric,
				Payload: nil,
			},
			expectError: true,
		},
		{
			name: "Payload is numeric",
			aggregate: QuestionAggregate{
				Type: TypeNumeric,
				Payload: &NumericQuestion{
					Correct: 1,
				},
			},
			expectError: false,
		},
		{
			name: "Payload is single",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Payload: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: 0,
				},
			},
			expectError: false,
		},
		{
			name: "Payload is multiple",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Payload: &MultipleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: []int{0},
				},
			},
			expectError: false,
		},
		{
			name: "Failed numeric validate",
			aggregate: QuestionAggregate{
				Type: TypeNumeric,
				Payload: &NumericQuestion{
					Correct: -1,
				},
			},
			expectError: true,
		},
		{
			name: "Failed single validate",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Payload: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: 23,
				},
			},
			expectError: true,
		},
		{
			name: "Failed multiple validate",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Payload: &MultipleChoiceQuestion{
					Options: nil,
					Correct: nil,
				},
			},
			expectError: true,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			err := tests[i].aggregate.Validate()
			if tests[i].expectError && err == nil {
				t.Errorf("expect error but got nil")
			}
			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
