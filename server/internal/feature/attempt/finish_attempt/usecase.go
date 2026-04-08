package finish_attempt

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type AttemptRepository interface {
	repository.Geter[*domain.Attempt]
	repository.Updater[*domain.Attempt]
}

type QuizRepository interface {
	repository.Geter[*domain.Quiz]
}

type UseCase struct {
	attemptRepo AttemptRepository
	quizRepo    QuizRepository
}

func New(attempt AttemptRepository, quiz QuizRepository) *UseCase {
	return &UseCase{
		attemptRepo: attempt,
		quizRepo:    quiz,
	}
}

func (u *UseCase) FinishAttempt(ctx context.Context, attemptID uuid.UUID, req Request) (Response, error) {
	attempt, err := u.attemptRepo.Get(ctx, attemptID)
	if err != nil {
		return Response{}, err
	}

	quiz, err := u.quizRepo.Get(ctx, attempt.QuizID)
	if err != nil {
		return Response{}, err
	}

	answerMap := make(map[uuid.UUID]domain.IAnswer)
	for _, answer := range attempt.SelectedAnswers {
		answerMap[answer.ID()] = answer
	}

	for _, req := range req.Answers {
		answer, ok := answerMap[req.QuestionID]
		if !ok {
			return Response{}, usecase.NewValidationError(fmt.Errorf("invalid question id"))
		}

		if err := req.WriteSelectedAnswer(answer); err != nil {
			return Response{}, err
		}
	}

	if err := attempt.Check(quiz); err != nil {
		return Response{}, err
	}

	attempt.Finish()

	if err := u.attemptRepo.Update(ctx, attempt); err != nil {
		return Response{}, err
	}

	return Response{
		Attempt: ResponseAttempt{
			ID:        attempt.ID,
			QuizID:    attempt.QuizID,
			UserID:    attempt.UserID,
			Score:     attempt.TotalScore,
			StartedAt: attempt.StartedAt,
			EndedAt:   attempt.EndedAt,
		},
	}, nil
}
