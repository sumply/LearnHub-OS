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
		OwnerID:   input.OwnerID,
		SubjectID: input.SubjectID,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	quiz.Prepare()

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
			Payload: u.createQuestionPayload(&input[i].Payload),
			Score:   1,
		}
		if err := question.Validate(); err != nil {
			return nil, err
		}
		content[i] = question
	}

	return content, nil
}

func (u *UseCase) createQuestionPayload(input *InputPayload) any {
	if input.Single != nil {
		return &domain.SingleChoiceQuestion{
			Options: input.Single.Options,
			Correct: input.Single.Correct,
		}
	}
	if input.Multiple != nil {
		return &domain.MultipleChoiceQuestion{
			Options: input.Multiple.Options,
			Correct: input.Multiple.Correct,
		}
	}
	if input.Numeric != nil {
		return &domain.NumericQuestion{
			Correct: input.Numeric.Correct,
		}
	}
	return nil
}
