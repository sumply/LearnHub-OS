package start_attempt

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"
)

type QuizRepository interface {
	repository.Geter[*domain.Quiz]
}

type AttemptRepository interface {
	repository.Saver[*domain.Attempt]
}

type UseCase struct {
	quizRepo    QuizRepository
	attemptRepo AttemptRepository
}

func New(quiz QuizRepository, attempt AttemptRepository) *UseCase {
	return &UseCase{
		quizRepo:    quiz,
		attemptRepo: attempt,
	}
}

func (u *UseCase) StartAttempt(ctx context.Context, req Request) (Response, error) {
	quiz, err := u.quizRepo.Get(ctx, req.QuizID)
	if err != nil {
		return Response{}, err
	}

	attempt, err := domain.NewAttempt(quiz, req.UserID)
	if err != nil {
		return Response{}, err
	}

	if err := u.attemptRepo.Save(ctx, attempt); err != nil {
		return Response{}, err
	}

	return Response{
		ID: attempt.ID,
	}, nil
}
