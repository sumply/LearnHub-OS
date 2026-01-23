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
	content, err := u.createDomainContent(input.Content)
	if err != nil {
		return Output{}, err
	}

	quiz := domain.Quiz{
		ID:        uuid.New(),
		Title:     input.Title,
		Summary:   input.Summary,
		OwnerID:   uuid.New(),
		SubjectID: input.SubjectID,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	err = u.postgres.CreateQuiz(ctx, &quiz)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: quiz.ID}, nil
}

func (u *UseCase) createDomainContent(input []InputContent) ([]domain.QuestionAggregate, error) {
	content := make([]domain.QuestionAggregate, len(input))
	for i := range content {
		question := domain.QuestionAggregate{
			ID:      uuid.New(),
			Type:    domain.QuestionType(input[i].Type),
			Text:    input[i].Text,
			Payload: input[i].Payload,
			Score:   1,
		}
		if err := question.Validate(); err != nil {
			return nil, err
		}
	}

	return content, nil
}
