package delete_quiz

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type Quiz interface {
	repository.Remover[domain.Quiz]
}

type UseCase struct {
	quiz Quiz
}

func New(quiz Quiz) *UseCase {
	return &UseCase{
		quiz: quiz,
	}
}

func (u *UseCase) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	err := u.quiz.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
