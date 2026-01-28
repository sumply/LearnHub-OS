package create_quiz

import (
	"context"
	"server/internal/domain"

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
	quizID := uuid.New()
	content, err := u.createDomainContent(quizID, input.Content)
	if err != nil {
		return Output{}, err
	}

	quiz, err := domain.NewQuiz(
		input.Title,
		input.Summary,
		input.OwnerID,
		input.SubjectID,
		content,
		input.MaxAttempts,
		input.Deadline,
	)
	if err != nil {
		return Output{}, err
	}
	quiz.ID = quizID

	err = u.postgres.CreateQuiz(ctx, &quiz)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: quiz.ID}, nil
}

func (u *UseCase) createDomainContent(quizID uuid.UUID, input []InputContent) ([]domain.QuestionAggregate, error) {
	content := make([]domain.QuestionAggregate, len(input))
	for i := range content {
		question := domain.QuestionAggregate{
			ID:      uuid.New(),
			QuizID:  quizID,
			Type:    domain.QuestionType(input[i].Type),
			Text:    input[i].Text,
			Details: u.createQuestionPayload(&input[i].Payload),
			Score:   1,
		}
		if err := question.Validate(); err != nil {
			return nil, err
		}
		content[i] = question
	}

	return content, nil
}

func (u *UseCase) createQuestionPayload(input *InputPayload) domain.QuestionDetails {
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
