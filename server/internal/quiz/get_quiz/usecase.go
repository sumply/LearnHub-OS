package get_quiz

import (
	"context"
	"server/internal/domain"
)

type Postgres interface {
	Quizzes(context.Context) ([]domain.Quiz, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) GetQuiz(ctx context.Context) (Output, error) {
	quizzes, err := u.postgres.Quizzes(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Quizzes: make([]OutputQuiz, len(quizzes)),
	}
	for i := range quizzes {
		output.Quizzes[i] = OutputQuiz{
			ID:         quizzes[i].ID,
			Title:      quizzes[i].Title,
			Summary:    quizzes[i].Summary,
			OwnerID:    quizzes[i].OwnerID,
			Content:    quizzes[i].Content,
			SubjectID:  quizzes[i].SubjectID,
			TotalScore: quizzes[i].TotalScore,
			CreatedAt:  quizzes[i].CreatedAt,
		}
	}

	return output, nil
}
