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
	log := s.loggerFromCreate(ctx, identity, name)
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
	log := s.loggerFromGet(ctx, identity)
	log.Debug("Called a get usecase method")

	domains, err := s.repo.Subject().GetAll(ctx)
	if err != nil {
		return nil, s.mapStorageError(err)
	}
	return domains, nil
}

func (s *RealSubject) loggerFromCreate(ctx context.Context, identity Identity, name string) logger.Logger {
	log := logger.FromCtx(ctx)
	field := s.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
		"name":     name,
	})
	return log.With(field)
}

func (s *RealSubject) loggerFromGet(ctx context.Context, identity Identity) logger.Logger {
	log := logger.FromCtx(ctx)
	field := s.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
	})
	return log.With(field)
}
