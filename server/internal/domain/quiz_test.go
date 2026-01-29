package domain

/*
func TestSingleChoiceQuestion(t *testing.T) {
	tests := []struct {
		Name        string
		question    SingleChoiceQuestion
		ExpectError bool
	}{
		{
			Name: "Options is empty",
			question: SingleChoiceQuestion{
				Options: nil,
				Correct: "",
			},
			ExpectError: true,
		},
		{
			Name: "Answer is incorrect",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: "invalid",
			},
			ExpectError: true,
		},
		{
			Name: "Answer is empty",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: "",
			}, ExpectError: true,
		},
		{
			Name: "Correct",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: "opt1",
			},
			ExpectError: false,
		},
		{
			Name: "Correct",
			question: SingleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: "opt2",
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
				Correct: []string{},
			},
			ExpectError: true,
		},
		{
			Name: "Answer is empty",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1"},
				Correct: nil,
			},
			ExpectError: true,
		},
		{
			Name: "Overflow",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []string{"", "", ""},
			},
			ExpectError: true,
		},
		{
			Name: "Answer is incorrect",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []string{"opt1", "invalid"},
			},
			ExpectError: true,
		},
		{
			Name: "Correct",
			question: MultipleChoiceQuestion{
				Options: []string{"opt1", "opt2"},
				Correct: []string{"opt1", "opt2"},
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
			expectError: false,
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
			name: "Invalid Details as numeric",
			aggregate: QuestionAggregate{
				Type: TypeNumeric,
				Details: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: "",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Details as single",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Details: &MultipleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: nil,
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Details as multiple",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Details: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: "",
				},
			},
			expectError: true,
		},
		{
			name: "Details is nil",
			aggregate: QuestionAggregate{
				Type:    TypeNumeric,
				Details: nil,
			},
			expectError: true,
		},
		{
			name: "Details is numeric",
			aggregate: QuestionAggregate{
				Type: TypeNumeric,
				Details: &NumericQuestion{
					Correct: 1,
				},
			},
			expectError: false,
		},
		{
			name: "Details is single",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Details: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: "opt1",
				},
			},
			expectError: false,
		},
		{
			name: "Details is multiple",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Details: &MultipleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: []string{"opt1"},
				},
			},
			expectError: false,
		},
		{
			name: "Failed single validate",
			aggregate: QuestionAggregate{
				Type: TypeSingleChoice,
				Details: &SingleChoiceQuestion{
					Options: []string{"opt1"},
					Correct: "",
				},
			},
			expectError: true,
		},
		{
			name: "Failed multiple validate",
			aggregate: QuestionAggregate{
				Type: TypeMultipleChoice,
				Details: &MultipleChoiceQuestion{
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

func TestNewQuiz(t *testing.T) {
	content := []QuestionAggregate{
		{
			Score: 1,
		},
		{
			Score: 5,
		},
		{
			Score: 2,
		},
	}
	tests := []struct {
		name        string
		param       newQuizParam
		expectError bool
	}{
		{
			name: "Content length is empty",
			param: newQuizParam{
				firstName:   "valid",
				lastName:    "valid",
				ownerID:     uuid.New(),
				subjectID:   uuid.New(),
				maxAttempts: 10,
				deadline:    nil,
				content:     nil,
			},
			expectError: true,
		},
		{
			name: "MaxAttempts is less 0",
			param: newQuizParam{
				firstName:   "valid",
				lastName:    "valid",
				ownerID:     uuid.New(),
				subjectID:   uuid.New(),
				maxAttempts: -1,
				deadline:    nil,
				content:     content,
			},
			expectError: true,
		},
		{
			name: "MaxAttempts is 0",
			param: newQuizParam{
				firstName:   "valid",
				lastName:    "valid",
				ownerID:     uuid.New(),
				subjectID:   uuid.New(),
				maxAttempts: 0,
				deadline:    nil,
				content:     content,
			},
			expectError: true,
		},
		{
			name: "Expired deadline",
			param: newQuizParam{
				firstName:   "valid",
				lastName:    "valid",
				ownerID:     uuid.New(),
				subjectID:   uuid.New(),
				maxAttempts: 1,
				deadline:    newExpiredDeadline(),
				content:     content,
			},
			expectError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewQuiz(
				test.param.firstName,
				test.param.lastName,
				test.param.ownerID,
				test.param.subjectID,
				test.param.content,
				test.param.maxAttempts,
				test.param.deadline,
			)
			if err != nil && !test.expectError {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && test.expectError {
				t.Error("expected error")
			}
		})
	}
}

func TestNewQuizCountTotalScore(t *testing.T) {
	content := []QuestionAggregate{
		{
			Score: 1,
		},
		{
			Score: 5,
		},
		{
			Score: 2,
		},
	}
	t.Run("Count total score", func(t *testing.T) {
		quiz, err := NewQuiz("title", "summary", uuid.New(), uuid.New(), content, 1, nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if quiz.TotalScore != 8 {
			t.Errorf("total score is invalid; expected: %d; got: %d", 8, quiz.TotalScore)
		}
	})
}

func newExpiredDeadline() *time.Time {
	t := time.Now().UTC().Add(-time.Hour)
	return &t
}

type newQuizParam struct {
	firstName   string
	lastName    string
	ownerID     uuid.UUID
	subjectID   uuid.UUID
	content     []QuestionAggregate
	maxAttempts int
	deadline    *time.Time
}
*/
