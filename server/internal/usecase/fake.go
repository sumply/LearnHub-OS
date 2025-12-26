package usecase

import (
	"context"
	"server/internal/domain"
)

type UserStub struct{}

func NewUserStub() *UserStub {
	return &UserStub{}
}

func (s *UserStub) Login(
	ctx context.Context,
	param UserLoginParam,
) (*domain.TokenPair, error) {
	return nil, ErrAccess
}

func (s *UserStub) Get(
	ctx context.Context,
	identity Identity,
) ([]*domain.User, error) {
	return nil, ErrAccess
}

func (s *UserStub) Create(
	ctx context.Context,
	identity Identity,
	param UserCreateParam,
) error {
	return ErrCollision
}

func (s *UserStub) GetMe(
	ctx context.Context,
	identity Identity,
) (*domain.User, error) {
	return nil, ErrNotFound
}

func (s *UserStub) GetByID(
	ctx context.Context,
	identity Identity,
	id domain.UserID,
) (*domain.User, error) {
	return nil, ErrNotFound
}

func (s *UserStub) Delete(
	ctx context.Context,
	identity Identity,
	id domain.UserID,
) error {
	return ErrNotFound
}

// =======================
// Speciality usecase stub
// =======================

type SpecialityStub struct{}

func NewSpecialityStub() *SpecialityStub {
	return &SpecialityStub{}
}

func (s *SpecialityStub) Create(
	ctx context.Context,
	identity Identity,
	name string,
) error {
	return ErrCollision
}

func (s *SpecialityStub) Get(
	ctx context.Context,
	identity Identity,
) ([]*domain.Speciality, error) {
	return nil, ErrAccess
}

// =======================
// Group usecase stub
// =======================

type GroupStub struct{}

func NewGroupStub() *GroupStub {
	return &GroupStub{}
}

func (s *GroupStub) Create(
	ctx context.Context,
	name string,
) error {
	return ErrCollision
}

func (s *GroupStub) Get(
	ctx context.Context,
) ([]*domain.Group, error) {
	return nil, ErrAccess
}

// =======================
// Subject usecase stub
// =======================

type SubjectStub struct{}

func NewSubjectStub() *SubjectStub {
	return &SubjectStub{}
}

func (s *SubjectStub) Create(
	ctx context.Context,
	identity Identity,
	param SubjectCreateParam,
) error {
	return ErrInvalidField
}

func (s *SubjectStub) Get(
	ctx context.Context,
	identity Identity,
) ([]*domain.Subject, error) {
	return nil, ErrAccess
}
