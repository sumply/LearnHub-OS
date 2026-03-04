package finish_attempt

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type Postgres interface {
	FinishedAttempt(context.Context, uuid.UUID) (dto.FinishedAttempt, error)
}

type Quiz interface {
	repository.Geter[domain.Quiz]
}

type Attempt interface {
	repository.Geter[domain.Attempt]
	repository.Updater[domain.Attempt]
}

type UseCase struct {
	postgres Postgres
	attempt  Attempt
	quiz     Quiz
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) FinishAttempt(ctx context.Context, attemptID uuid.UUID, input *Input) (Output, error) {
	attempt, err := u.attempt.Get(ctx, attemptID)
	if err != nil {
		return Output{}, err
	}

	quiz, err := u.quiz.Get(ctx, attempt.QuizID)
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

	err = u.attempt.Update(ctx, attempt)
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
		a := &answers[i]
		answersMap[a.QuestionID] = a
	}

	for i := range input {
		answer, ok := answersMap[input[i].QuestionID]
		if !ok {
			return fmt.Errorf("unknown questionID (id=%s)", input[i].QuestionID)
		}

		answer.Answer = input[i].Answer
	}

	return nil
}
