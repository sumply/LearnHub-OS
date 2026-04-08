package list_quizzes

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"
)

type QuizRepository interface {
	repository.AllGeter[*domain.Quiz]
}

type UseCase struct {
	quizRepo QuizRepository
}

func New(quiz QuizRepository) *UseCase {
	return &UseCase{
		quizRepo: quiz,
	}
}

func (u *UseCase) ListQuizzes(ctx context.Context) (Response, error) {
	quizzes, err := u.quizRepo.GetAll(ctx)
	if err != nil {
		return Response{}, err
	}

	respQuizzes := make([]ResponseQuiz, 0, len(quizzes))
	for _, quiz := range quizzes {
		respQuizzes = append(respQuizzes, NewResponseQuiz(quiz))
	}

	return Response{
		Quizzes: respQuizzes,
	}, nil
}
