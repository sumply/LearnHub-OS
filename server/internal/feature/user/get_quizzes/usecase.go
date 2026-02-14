package get_quizzes

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetByStudent(context.Context, uuid.UUID) (domain.Group, error)
}

type QuizItemRepository interface {
	ListByGroup(context.Context, uuid.UUID) ([]dto.QuizItem, error)
	ListByOwner(context.Context, uuid.UUID) ([]dto.QuizItem, error)
}

type Repository struct {
	Group    GroupRepository
	QuizItem QuizItemRepository
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) GetQuizzes(ctx context.Context, userID uuid.UUID, opts Options) (Output, error) {
	switch opts.Scope {
	case ScopeOwner:
		return u.handleScopeOwner(ctx, userID)
	case ScopeAvailable:
		return u.handleScopeAvailable(ctx, userID)
	}
	return Output{}, nil
}

func (u *UseCase) handleScopeOwner(ctx context.Context, userID uuid.UUID) (Output, error) {
	quizzes, err := u.repository.QuizItem.ListByOwner(ctx, userID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Quizzes: quizzes,
	}, nil
}

func (u *UseCase) handleScopeAvailable(ctx context.Context, userID uuid.UUID) (Output, error) {
	group, err := u.repository.Group.GetByStudent(ctx, userID)
	if err != nil {
		return Output{}, err
	}

	quizzes, err := u.repository.QuizItem.ListByGroup(ctx, group.ID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Quizzes: quizzes,
	}, nil
}
