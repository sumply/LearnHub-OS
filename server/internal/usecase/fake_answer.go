package usecase

import "context"

type FakeQuizResult struct{}

func NewFakeQuizResult() *FakeQuizResult {
	return &FakeQuizResult{}
}

func (f *FakeQuizResult) Create(context.Context, Identity, AnswerCreateParam) error {
	return nil
}
func (f *FakeQuizResult) Get(context.Context, Identity) ([]AnswerDomain, error) {
	d := AnswerDomain{
		ID:         0,
		TotalScore: 0,
		Score:      0,
		Completed:  false,
		Quiz: QuizDomain{
			ID:      0,
			Name:    "Name",
			Summary: "Summary",
			Questions: []QuizQuestionDomain{
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
			},
		},
	}
	res := make([]AnswerDomain, 10)
	for i := range 10 {
		res[i] = d
		if i&1 == 0 {
			res[i].Completed = true
		}
	}
	return res, nil
}
func (f *FakeQuizResult) GetByID(context.Context, Identity, ID) (AnswerDomain, error) {
	d := AnswerDomain{
		ID:         0,
		TotalScore: 0,
		Score:      0,
		Completed:  false,
		Quiz: QuizDomain{
			ID:      0,
			Name:    "Name",
			Summary: "Summary",
			Questions: []QuizQuestionDomain{
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
				{
					ID:   0,
					Name: "Name",
					Answers: []QuizOptionsDomain{
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: false,
						},
						{
							ID:        0,
							Text:      "Text",
							IsCorrect: true,
						},
					},
				},
			},
		},
	}
	return d, nil
}
