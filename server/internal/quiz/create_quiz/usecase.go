package create_quiz

import (
	"context"
	"fmt"
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

	content, err := u.createQuizContent(quizID, input)
	if err != nil {
		return Output{}, err
	}

	quiz, err := domain.NewQuiz(
		quizID,
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

func (u *UseCase) createQuestionDetails(content *InputContent) (domain.QuestionDetails, error) {
	switch domain.QuestionType(content.Type) {
	case domain.TypeSingleChoice:
		single, err := domain.NewSingleChoiceQuestion(
			content.Payload.Single.Options,
			content.Payload.Single.Correct,
		)
		if err != nil {
			return nil, err
		}
		return &single, nil
	case domain.TypeMultipleChoice:
		multiple, err := domain.NewMultipleChoiceQuestion(
			content.Payload.Multiple.Options,
			content.Payload.Multiple.Correct,
		)
		if err != nil {
			return nil, err
		}
		return &multiple, nil
	case domain.TypeNumeric:
		numeric := domain.NewNumericQuestion(
			content.Payload.Numeric.Correct,
		)
		return &numeric, nil
	default:
		return nil, fmt.Errorf("invalid type")
	}
}

func (u *UseCase) createQuizContent(quizID uuid.UUID, input *Input) ([]domain.QuestionAggregate, error) {
	content := make([]domain.QuestionAggregate, len(input.Content))
	for i := range content {
		details, err := u.createQuestionDetails(&input.Content[i])
		if err != nil {
			return nil, err
		}

		question, err := domain.NewQuestionAggregate(quizID, input.Content[i].Text, details, input.Content[i].Score)
		if err != nil {
			return nil, err
		}

		content[i] = question
	}

	return content, nil
}
