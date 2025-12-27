package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
)

type RealSpeciality struct {
	usecase
	repo *repository.Repository
}

func NewRealSpeciality(repo *repository.Repository) *RealSpeciality {
	return &RealSpeciality{
		repo: repo,
	}
}

func (s *RealSpeciality) Create(ctx context.Context, identity Identity, name string) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(name),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		log.Warn("User role is less than admin")
		return ErrAccess
	}

	spec, err := domain.NewSpeciality(name)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	ctx = logger.WithLoggerCtx(ctx, log)
	err = s.repo.Speciality().Save(ctx, spec)
	if err != nil {
		log.Warn(err.Error())
		return s.mapStorageError(err)
	}

	return nil
}

func (s *RealSpeciality) Get(ctx context.Context, identity Identity) ([]*domain.Speciality, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase method")

	ctx = logger.WithLoggerCtx(ctx, log)
	domains, err := s.repo.Speciality().GetAll(ctx)
	if err != nil {
		log.Warn(err.Error())
		return nil, s.mapStorageError(err)
	}

	return domains, nil
}
