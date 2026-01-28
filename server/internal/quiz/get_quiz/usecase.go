package get_quiz

import (
	"context"
	"server/internal/dto"
)

type Postgres interface {
	QuizItems(context.Context) ([]dto.QuizItem, error)
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
	quizzes, err := u.postgres.QuizItems(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Quizzes: quizzes,
	}

	return output, nil
}
