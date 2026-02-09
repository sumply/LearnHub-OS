package find_quiz

import (
	"context"
	"server/internal/dto"
	"server/internal/repository/filter"
)

type Repository interface {
	FindQuizItems(context.Context, *filter.Quiz) ([]dto.QuizItem, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) FindForStudent(ctx context.Context, request *Request) (Response, error) {
	quizzes, err := u.repository.FindQuizItems(ctx, &filter.Quiz{
		Attempt: &filter.Attempt{
			UserID:      &request.UserID,
			IsCompleted: request.IsCompleted,
		},
	})
	if err != nil {
		return Response{}, err
	}
	return Response{
		Quizzes: quizzes,
	}, nil
}
