package get_attempt

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type AttemptRepository interface {
	repository.Geter[*domain.Attempt]
}

type UseCase struct {
	attemptRepo AttemptRepository
}

func New(attempt AttemptRepository) *UseCase {
	return &UseCase{
		attemptRepo: attempt,
	}
}

func (u *UseCase) GetAttempt(ctx context.Context, attemptID uuid.UUID) (Response, error) {
	attempt, err := u.attemptRepo.Get(ctx, attemptID)
	if err != nil {
		return Response{}, err
	}

	respAnswers := make([]ResponseAnswer, 0, len(attempt.SelectedAnswers))
	for _, answer := range attempt.SelectedAnswers {
		var newRespAnswer ResponseAnswer
		err := answer.Accept(&newRespAnswer)
		if err != nil {
			return Response{}, err
		}
		respAnswers = append(respAnswers, newRespAnswer)
	}

	return Response{
		Attempt: ResponseAttempt{
			ID:        attempt.ID,
			QuizID:    attempt.QuizID,
			UserID:    attempt.UserID,
			Score:     attempt.TotalScore,
			StartedAt: attempt.StartedAt,
			Answers:   respAnswers,
			EndedAt:   attempt.EndedAt,
		},
	}, nil
}
