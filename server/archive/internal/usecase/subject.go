package usecase

import (
	"context"
	"server/archive/internal/domain"
	"server/archive/internal/dto"
	"server/archive/internal/logger"
	"server/archive/internal/repository"
)

type RealSubject struct {
	usecase
	repo *repository.Repository
}

func NewRealSubject(repo *repository.Repository) *RealSubject {
	return &RealSubject{
		repo: repo,
	}
}

func (s *RealSubject) Create(ctx context.Context, identity *dto.Identity, req *dto.SubjectCreateReq) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(req),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}

	domain, err := domain.NewSubject(req.Name)
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

func (s *RealSubject) Get(ctx context.Context, identity *dto.Identity) ([]*domain.Subject, error) {
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
