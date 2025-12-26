package repository

import (
	"context"
	"errors"
	"server/internal/domain"
)

var (
	ErrCollision = errors.New("collision")
)

type GroupFilter struct {
	ID          domain.GroupID
	WithCurator bool
}

type User interface {
	Save(context.Context, *domain.User) error
	GetByID(context.Context, domain.UserID) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
}

type Speciality interface {
	Save(context.Context, *domain.Speciality) error
	GetAll(context.Context) ([]*domain.Speciality, error)
}

type Subject interface {
	Save(context.Context, *domain.Subject) error
	GetAll(context.Context) ([]*domain.Subject, error)
}

type Group interface {
	Save(context.Context, *domain.Group) error
	GetAll(context.Context) ([]*domain.Group, error)
	AddStudent(context.Context, domain.GroupID, []domain.UserID) error
	RemoveStudent(context.Context, domain.GroupID, []domain.UserID) error
	Find(context.Context, GroupFilter) ([]*domain.Group, error)
}

type Repository interface {
	User() User
	Subject() Subject
	Group() Group
	Speciality() Speciality
}
