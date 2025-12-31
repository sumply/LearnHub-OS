package usecase

import (
	"context"
	"errors"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/dto"
)

var (
	ErrAccess       = errors.New("access permission")
	ErrNotFound     = errors.New("not found")
	ErrCollision    = errors.New("collision")
	ErrInvalidField = errors.New("invalid field")
)

type QuizInterface interface {
	Create(context.Context, *dto.Identity, *dto.QuizCreateReq) error
	Get(context.Context, *dto.Identity) ([]*domain.Quiz, error)
}

type UserInterface interface {
	Login(context.Context, *dto.LoginReq) (*domain.TokenPair, error)
	Get(context.Context, *dto.Identity) ([]*domain.User, error)
	Create(context.Context, *dto.Identity, *dto.UserCreateReq) error
	GetMe(context.Context, *dto.Identity) (*domain.User, error)
	GetByID(context.Context, *dto.Identity, common.ID) (*domain.User, error)
}

type GroupInterface interface {
	Create(context.Context, *dto.Identity, *dto.GroupCreateReq) error
	Get(ctx context.Context) ([]*domain.Group, error)
}

type SubjectInterface interface {
	Create(context.Context, *dto.Identity, *dto.SubjectCreateReq) error
	Get(context.Context, *dto.Identity) ([]*domain.Subject, error)
}
