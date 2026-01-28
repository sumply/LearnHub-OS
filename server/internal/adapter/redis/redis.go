package redis

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Redis struct {
	attempts map[uuid.UUID]*domain.Attempt
}

func (r *Redis) SaveAttempt(ctx context.Context, attempt *domain.Attempt) error {
	r.attempts[attempt.ID] = attempt
	return nil
}
