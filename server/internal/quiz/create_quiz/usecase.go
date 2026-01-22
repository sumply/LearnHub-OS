package create_quiz

import (
	"context"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Postgres interface {
	CreateQuiz(context.Context, *domain.Quiz) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) CreateQuiz(ctx context.Context, input *Input) (Output, error) {
	quiz := domain.Quiz{
		ID:        uuid.New(),
		Title:     input.Title,
		Summary:   input.Summary,
		OwnerID:   uuid.New(),
		SubjectID: input.SubjectID,
		Content:   input.Content,
		CreatedAt: time.Now().UTC(),
	}

	err := u.postgres.CreateQuiz(ctx, &quiz)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: quiz.ID}, nil
}
