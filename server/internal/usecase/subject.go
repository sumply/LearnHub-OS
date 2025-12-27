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

func (s *RealSubject) Create(ctx context.Context, identity Identity, param SubjectCreateParam) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(param),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}

	domain, err := domain.NewSubject(param.Name, param.SpecialityIDs)
	if err != nil {
		return err
	}

	ctx = logger.WithLoggerCtx(ctx, log)
	err = s.repo.Subject().Save(ctx, domain)
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

	ctx = logger.WithLoggerCtx(ctx, log)
	domains, err := s.repo.Subject().GetAll(ctx)
	if err != nil {
		return nil, s.mapStorageError(err)
	}
	return domains, nil
}
