package list_attempts

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"
)

type AttemptRepository interface {
	repository.AllGeter[*domain.Attempt]
}

type UseCase struct {
	attemptRepo AttemptRepository
}

func New(attempt AttemptRepository) *UseCase {
	return &UseCase{
		attemptRepo: attempt,
	}
}

func (u *UseCase) ListAttempts(ctx context.Context) (Response, error) {
	attempts, err := u.attemptRepo.GetAll(ctx)
	if err != nil {
		return Response{}, err
	}

	attemptsResp := make([]ResponseAttempt, 0, len(attempts))
	for _, attempt := range attempts {
		attemptsResp = append(attemptsResp, ResponseAttempt{
			ID:        attempt.ID,
			UserID:    attempt.UserID,
			QuizID:    attempt.QuizID,
			Score:     attempt.TotalScore,
			StartedAt: attempt.StartedAt,
			EndedAt:   attempt.EndedAt,
		})
	}

	return Response{
		Attempts: attemptsResp,
	}, nil
}
