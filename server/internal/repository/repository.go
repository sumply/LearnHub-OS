package repository

import (
	"context"
	"errors"
	"server/internal/domain"
)

var (
	ErrCollision  = errors.New("collision")
	ErrNotFound   = errors.New("not found")
	ErrDependensy = errors.New("dependency")
	ErrInvalid    = errors.New("invalid")
)

type GroupFilter struct {
	ID          domain.GroupID
	WithCurator bool
}

type User interface {
	Save(context.Context, *domain.User) error
	GetByID(context.Context, domain.UserID) (*domain.User, error)
	GetByLogin(context.Context, domain.Login) (*domain.User, error)
	GetAll(context.Context) ([]*domain.User, error)
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
}

type Repository struct {
	user       User
	subject    Subject
	group      Group
	speciality Speciality
}

func New(u User, s Subject, g Group, sp Speciality) *Repository {
	if u == nil || s == nil || g == nil || sp == nil {
		panic("Repository params is nil")
	}
	return &Repository{
		user:       u,
		subject:    s,
		group:      g,
		speciality: sp,
	}
}

func (r *Repository) User() User {
	return r.user
}

func (r *Repository) Subject() Subject {
	return r.subject
}

func (r *Repository) Group() Group {
	return r.group
}

func (r *Repository) Speciality() Speciality {
	return r.speciality
}
