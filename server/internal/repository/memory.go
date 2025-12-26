package repository

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/logger"
)

var userStorage = make(map[domain.UserID]*domain.User)

type UserMemory struct{}

func NewUserMemory() User {
	return &UserMemory{}
}

func (u *UserMemory) Save(ctx context.Context, user *domain.User) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(user),
	)
	log.Debug("Called a save repository method")
	if _, ok := userStorage[user.ID]; ok {
		return fmt.Errorf("%s: user_id=%d", ErrCollision, user.ID)
	}
	userStorage[user.ID] = user
	return nil
}

func (u *UserMemory) GetByID(context.Context, domain.UserID) (*domain.User, error) {
	return nil, nil
}
func (u *UserMemory) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	return nil, nil
}
