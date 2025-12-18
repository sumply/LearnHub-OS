package usecase

import "context"

type StubQuiz struct{}

func NewStubQuiz() *StubQuiz {
	return &StubQuiz{}
}

func (f *StubQuiz) Create(ctx context.Context, auth Identity, param QuizCreateParam) error {
	return nil
}

func (f *StubQuiz) Get(ctx context.Context, auth Identity) ([]QuizDomain, error) {
	answer := QuizOptionsDomain{
		ID:        0,
		Text:      "Text",
		IsCorrect: true,
	}
	question := QuizQuestionDomain{
		ID:   0,
		Name: "Title",
		Answers: []QuizOptionsDomain{
			answer,
			answer,
			answer,
			answer,
		},
	}
	quiz := QuizDomain{
		ID:      0,
		Name:    "Name",
		Summary: "Summary",
		Questions: []QuizQuestionDomain{
			question,
			question,
			question,
			question,
			question,
			question,
		},
	}
	slice := []QuizDomain{
		quiz,
		quiz,
		quiz,
		quiz,
		quiz,
	}
	return slice, nil
}

func (f *StubQuiz) GetByID(ctx context.Context, auth Identity, id ID) (QuizDomain, error) {
	answer := QuizOptionsDomain{
		ID:        0,
		Text:      "Text",
		IsCorrect: true,
	}
	question := QuizQuestionDomain{
		ID:   0,
		Name: "Title",
		Answers: []QuizOptionsDomain{
			answer,
			answer,
			answer,
			answer,
		},
	}
	quiz := QuizDomain{
		ID:      0,
		Name:    "Name",
		Summary: "Summary",
		Questions: []QuizQuestionDomain{
			question,
			question,
			question,
			question,
			question,
			question,
		},
	}
	return quiz, nil
}
