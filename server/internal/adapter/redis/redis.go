package redis

import (
	"context"
	"fmt"
	"server/internal/domain"
)

type Redis struct{}

func (r *Redis) SaveAttempt(ctx context.Context, attempt *domain.Attempt) error {
	fmt.Println(attempt)
	return nil
}
