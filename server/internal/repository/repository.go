package repository

import (
	"context"
	"errors"
	"server/internal/common"
	"server/internal/domain"
)

var (
	ErrCollision  = errors.New("collision")
	ErrNotFound   = errors.New("not found")
	ErrDependensy = errors.New("dependency")
	ErrInvalid    = errors.New("invalid")
)

type GroupFilter struct {
	ID          common.ID
	WithCurator bool
}

type UserInterface interface {
	Save(context.Context, *domain.User) error
	GetByID(context.Context, common.ID) (*domain.User, error)
	GetByLogin(context.Context, domain.Login) (*domain.User, error)
	GetAll(context.Context) ([]*domain.User, error)
}

type SubjectInterface interface {
	Save(context.Context, *domain.Subject) error
	GetAll(context.Context) ([]*domain.Subject, error)
}

type GroupInterface interface {
	Save(context.Context, *domain.Group) error
	GetAll(context.Context) ([]*domain.Group, error)
	AddStudent(context.Context, common.ID, []common.ID) error
	RemoveStudent(context.Context, common.ID, []common.ID) error
}

type QuizInterface interface {
	Save(context.Context, *domain.Quiz) error
}

type Repository struct {
	user    UserInterface
	subject SubjectInterface
	group   GroupInterface
	quiz    QuizInterface
}

func New(u UserInterface, s SubjectInterface, g GroupInterface, q QuizInterface) *Repository {
	if u == nil || s == nil || g == nil {
		panic("Repository params is nil")
	}
	return &Repository{
		user:    u,
		subject: s,
		group:   g,
		quiz:    q,
	}
}

func (r *Repository) User() UserInterface {
	return r.user
}

func (r *Repository) Subject() SubjectInterface {
	return r.subject
}

func (r *Repository) Group() GroupInterface {
	return r.group
}

func (r *Repository) Quiz() QuizInterface {
	return r.quiz
}
