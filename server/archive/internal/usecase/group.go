package usecase

import (
	"context"
	"server/archive/internal/common"
	"server/archive/internal/domain"
	"server/archive/internal/dto"
	"server/archive/internal/logger"
	"server/archive/internal/repository"
)

type GroupReal struct {
	repo *repository.Repository
}

func NewGroupReal(repo *repository.Repository) *GroupReal {
	return &GroupReal{
		repo: repo,
	}
}

func (g *GroupReal) Create(ctx context.Context, identity *dto.Identity, req *dto.GroupCreateReq) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(req),
	)
	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		log.Warn("User role is less than admin")
		return ErrAccess
	}

	group, err := domain.NewGroup(req.Name, req.CuratorID)
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

func (g *GroupReal) AddStudents(
	ctx context.Context,
	identity *dto.Identity,
	groupID common.ID,
	req *dto.GroupAddStudentsReq,
) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(req),
	)
	log.Debug("Called a addStudents groupReal method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		log.Warn("A user role less a admin")
		return ErrAccess
	}

	err := g.repo.Group().AddStudent(ctx, groupID, req.StudentIDs)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func (g *GroupReal) GetByID(ctx context.Context, identity *dto.Identity, groupID common.ID) (*domain.Group, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.NewTracedField("groupID", groupID),
	)
	log.Debug("Called a getByID group usecase")

	group, err := g.repo.Group().GetByID(logger.WithLoggerCtx(ctx, log), groupID)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}
	return group, nil
}

func (g *GroupReal) DeleteByID(ctx context.Context, identity *dto.Identity, groupID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("groupID", groupID),
	)
	log.Warn("Called a deleteByID group usecase")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}

	err := g.repo.Group().DeleteByID(ctx, groupID)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (g *GroupReal) DeleteStudentByID(ctx context.Context, identity *dto.Identity, groupID, studentID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.NewTracedField("studentID", studentID),
	)
	log.Debug("Called a deleteStudentByID group usecase")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}

	err := g.repo.Group().RemoveStudent(ctx, groupID, studentID)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	return nil
}
