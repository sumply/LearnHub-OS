package repository

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/logger"
	"time"
)

var credStorage = make(map[domain.Login]*domain.Credential)
var userStorage = make(map[domain.UserID]*domain.User)
var nextID = 1

func init() {
	credStorage["vlad"] = &domain.Credential{
		Login:     "vlad",
		PwdHashed: "verysecret",
		Email:     "vlad@gmail.com",
	}
	userStorage[0] = &domain.User{
		ID:         0,
		FirstName:  "Vladislav",
		LastName:   "Yanushkevich",
		MiddleName: "Vitalevich",
		Role:       domain.UserRoot,
		Credential: credStorage["vlad"],
		CreatedAt:  time.Now().UTC(),
	}
}

type UserMemory struct{}

func NewUserMemory() User {
	return &UserMemory{}
}

func (u *UserMemory) Save(ctx context.Context, user *domain.User) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(user),
	)
	log.Debug("Called a save repository method")
	user.ID = domain.UserID(nextID)
	if _, ok := credStorage[user.Credential.Login]; ok {
		return fmt.Errorf("%w: login=%s", ErrCollision, user.Credential.Login)
	}
	credStorage[user.Credential.Login] = user.Credential
	if _, ok := userStorage[user.ID]; ok {
		return fmt.Errorf("%w: user_id=%d", ErrCollision, user.ID)
	}
	userStorage[user.ID] = user
	nextID++
	return nil
}

func (u *UserMemory) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(id),
	)
	log.Debug("Called a save repository method")
	user, ok := userStorage[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}
func (u *UserMemory) GetByLogin(ctx context.Context, login domain.Login) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(login),
	)
	log.Debug("Called a getByLogin repository method")
	for _, user := range userStorage {
		if user.Credential.Login == login {
			return user, nil
		}
	}
	return nil, ErrNotFound
}
func (u *UserMemory) GetAll(ctx context.Context) ([]*domain.User, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(userStorage) == 0 {
		return nil, ErrNotFound
	}

	users := make([]*domain.User, len(userStorage))
	for i, user := range userStorage {
		u := *user
		u.Credential = nil
		users[i] = &u
	}

	return users, nil
}
