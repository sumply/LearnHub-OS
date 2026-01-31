package finish_attempt

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Postgres interface {
	DomainQuiz(context.Context, uuid.UUID) (domain.Quiz, error)
	DomainAttempt(context.Context, uuid.UUID) (domain.Attempt, error)
	UpdateAttempt(context.Context, *domain.Attempt) error
	FinishedAttempt(context.Context, uuid.UUID) (dto.FinishedAttempt, error)
}
type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) FinishAttempt(ctx context.Context, attemptID uuid.UUID, input *Input) (Output, error) {
	attempt, err := u.postgres.DomainAttempt(ctx, attemptID)
	if err != nil {
		return Output{}, err
	}

	quiz, err := u.postgres.DomainQuiz(ctx, attempt.QuizID)
	if err != nil {
		return Output{}, err
	}

	err = u.applyAnswers(attempt.Answers, input.Answers)
	if err != nil {
		return Output{}, err
	}

	err = quiz.CheckAttempt(&attempt)
	if err != nil {
		return Output{}, err
	}

	attempt.Finish()

	err = u.postgres.UpdateAttempt(ctx, &attempt)
	if err != nil {
		return Output{}, err
	}

	attemptDTO, err := u.postgres.FinishedAttempt(ctx, attempt.ID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		FinishedAttempt: attemptDTO,
	}, nil
}

func (u *UseCase) applyAnswers(answers []domain.Answer, input []InputAnswer) error {
	answersMap := make(map[uuid.UUID]*domain.Answer)
	for i := range answers {
		answersMap[answers[i].QuestionID] = &answers[i]
	}

	for i := range input {
		answer, ok := answersMap[input[i].QuestionID]
		if !ok {
			return fmt.Errorf("unknown questionID (id=%d)", input[i].QuestionID)
		}

		answer.Answer = input[i].Answer
	}

	return nil
}
