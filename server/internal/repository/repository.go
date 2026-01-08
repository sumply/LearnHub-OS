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
	ErrDependence = errors.New("dependence")
	ErrInvalid    = errors.New("invalid")
)

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

type GroupFilter struct {
	StudentID *common.ID
}

type GroupInterface interface {
	Save(context.Context, *domain.Group) error
	GetAll(context.Context) ([]*domain.Group, error)
	GetByID(context.Context, common.ID) (*domain.Group, error)
	AddStudent(context.Context, common.ID, []common.ID) error
	RemoveStudent(context.Context, common.ID, []common.ID) error
}

type QuizFilter struct {
	OwnerID *common.ID
	Group   *GroupFilter
}

type QuizInterface interface {
	Save(context.Context, *domain.Quiz) error
	GetAll(context.Context) ([]*domain.Quiz, error)
	GetByID(context.Context, common.ID) (*domain.Quiz, error)
	GetWithFilter(context.Context, *QuizFilter) ([]*domain.Quiz, error)
	Delete(context.Context, common.ID) error
}

type ProgressFilter struct {
	UserID *common.ID
	Quiz   *QuizFilter
}

type ProgressInterface interface {
	UpdateAnswer(context.Context, *domain.Answer) error
	GetByID(context.Context, common.ID) (*domain.QuizProgress, error)
	Get(context.Context, *ProgressFilter) ([]*domain.QuizProgress, error)
	Update(*domain.QuizProgress) error
}

type Repository struct {
	user     UserInterface
	subject  SubjectInterface
	group    GroupInterface
	quiz     QuizInterface
	progress ProgressInterface
}

func New(u UserInterface, s SubjectInterface, g GroupInterface, q QuizInterface, p ProgressInterface) *Repository {
	if u == nil || s == nil || g == nil {
		panic("Repository params is nil")
	}
	return &Repository{
		user:     u,
		subject:  s,
		group:    g,
		quiz:     q,
		progress: p,
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

func (r *Repository) Progress() ProgressInterface {
	return r.progress
}
