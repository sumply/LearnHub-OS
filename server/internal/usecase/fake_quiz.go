package usecase

import "context"

type FakeQuiz struct{}

func NewFakeQuiz() *FakeQuiz {
	return &FakeQuiz{}
}

func (f *FakeQuiz) Create(ctx context.Context, auth Identity, param QuizCreateParam) error {
	return nil
}

func (f *FakeQuiz) Get(ctx context.Context, auth Identity) ([]QuizDomain, error) {
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

func (f *FakeQuiz) GetByID(ctx context.Context, auth Identity, id ID) (QuizDomain, error) {
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

type FakeQuizResult struct{}

func NewFakeQuizResult() *FakeQuizResult {
	return &FakeQuizResult{}
}

func (f *FakeQuizResult) Create(context.Context, Identity, []QuizResultCreateParam) error {
	return nil
}
func (f *FakeQuizResult) Get(context.Context, Identity) error {
	return nil
}
func (f *FakeQuizResult) GetByID(context.Context, Identity, ID) error {
	return nil
}
