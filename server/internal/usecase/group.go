package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
)

type GroupReal struct {
	repo *repository.Repository
}

func NewGroupReal(repo *repository.Repository) *GroupReal {
	return &GroupReal{
		repo: repo,
	}
}

func (g *GroupReal) Create(ctx context.Context, identity Identity, param GroupCreateParam) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(param),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		log.Warn("User role is less than admin")
		return ErrAccess
	}

	group, err := domain.NewGroup(param.Name, param.CuratorID, param.SpecialityID)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	ctx = logger.WithLoggerCtx(ctx, log)
	err = g.repo.Group().Save(ctx, group)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (g *GroupReal) Get(ctx context.Context) ([]*domain.Group, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a get usecase method")

	groups, err := g.repo.Group().GetAll(ctx)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}

	return groups, nil
}
