package repository

import (
	"context"
	"errors"
	"server/internal/domain"
)

// =======================
// User stub
// =======================

type UserStub struct{}

func NewUserStub() *UserStub {
	return &UserStub{}
}

func (s *UserStub) Save(ctx context.Context, u *domain.User) error {
	return nil
}

func (s *UserStub) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return nil, errors.New("user not found")
}

func (s *UserStub) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	return nil, errors.New("user not found")
}

// =======================
// Speciality stub
// =======================

type SpecialityStub struct{}

func NewSpecialityStub() *SpecialityStub {
	return &SpecialityStub{}
}

func (s *SpecialityStub) Save(ctx context.Context, sp *domain.Speciality) error {
	return nil
}

func (s *SpecialityStub) GetAll(ctx context.Context) ([]*domain.Speciality, error) {
	return []*domain.Speciality{}, nil
}

// =======================
// Subject stub
// =======================

type SubjectStub struct{}

func NewSubjectStub() *SubjectStub {
	return &SubjectStub{}
}

func (s *SubjectStub) Save(ctx context.Context, subj *domain.Subject) error {
	return nil
}

func (s *SubjectStub) GetAll(ctx context.Context) ([]*domain.Subject, error) {
	return []*domain.Subject{}, nil
}

// =======================
// Group stub
// =======================

type GroupStub struct{}

func NewGroupStub() *GroupStub {
	return &GroupStub{}
}

func (s *GroupStub) Save(ctx context.Context, g *domain.Group) error {
	return nil
}

func (s *GroupStub) GetAll(ctx context.Context) ([]*domain.Group, error) {
	return []*domain.Group{}, nil
}

func (s *GroupStub) AddStudent(
	ctx context.Context,
	groupID domain.GroupID,
	userIDs []domain.UserID,
) error {
	return nil
}

func (s *GroupStub) RemoveStudent(
	ctx context.Context,
	groupID domain.GroupID,
	userIDs []domain.UserID,
) error {
	return nil
}

func (s *GroupStub) Find(ctx context.Context, filter GroupFilter) ([]*domain.Group, error) {
	return []*domain.Group{}, nil
}
