package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/validator"
)

type RealSubject struct {
	usecase
	val  validator.Subject
	repo repository.Repository
}

func NewRealSubject() *RealSubject {
	return &RealSubject{}
}

func (s *RealSubject) Create(ctx context.Context, identity Identity, name string) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(name),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}
	name = s.val.Trim(name)
	if !s.val.ValidName(name) {
		return ErrInvalidField
	}

	err := s.repo.Subject().Save(ctx, nil)
	if err != nil {
		return s.mapStorageError(err)
	}
	return nil
}

func (s *RealSubject) Get(ctx context.Context, identity Identity) ([]*domain.Subject, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase method")

	domains, err := s.repo.Subject().GetAll(ctx)
	if err != nil {
		return nil, s.mapStorageError(err)
	}
	return domains, nil
}
