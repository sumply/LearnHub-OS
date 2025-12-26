package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/generator"
)

func NewRealUser(gen generator.Generator, repo repository.Repository) *RealUser {
	return &RealUser{
		gen:  gen,
		repo: repo,
	}
}

type RealUser struct {
	usecase
	repo repository.Repository
	gen  generator.Generator
}

func (u *RealUser) Login(ctx context.Context, param UserLoginParam) (*domain.TokenPair, error) {
	log := logger.FromCtx(ctx).
		With(
			logger.TraceFieldFromAny(param),
		)
	log.Debug("Called a login usecase method")

	user, err := u.repo.User().GetByLogin(ctx, param.Login)
	if err != nil {
		log.Warn(err.Error())
		return nil, u.mapStorageError(err)
	}

	ok := user.Credential.Authorization(param.Login, param.Password)
	if !ok {
		return nil, ErrAccess
	}

	access, refresh := u.gen.GenJWTTokens(uint64(user.ID), string(user.Role))
	log.Debug("Generated jwt tokens")

	return &domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (u *RealUser) Get(ctx context.Context, identity Identity) ([]*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase method")
	return nil, nil
}

func (u *RealUser) Create(ctx context.Context, identity Identity, param UserCreateParam) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(param),
	)

	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) || !identity.Role.IsHigher(param.Role) {
		log.Warn("Access permission")
		return ErrAccess
	}

	login := u.gen.GenLogin()
	pwd := u.gen.GenPassword()

	user, err := domain.NewUser(
		login,
		pwd,
		param.Email,
		param.FirstName,
		param.LastName,
		param.MiddleName,
		domain.UserRole(param.Role),
	)
	if err != nil {
		return err
	}

	err = u.repo.User().Save(ctx, user)
	if err != nil {
		return u.mapStorageError(err)
	}

	return nil
}

func (u *RealUser) GetMe(ctx context.Context, identity Identity) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a getMe usecase method")

	user, err := u.repo.User().GetByID(ctx, domain.UserID(identity.ID))
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	return user, nil
}

func (u *RealUser) GetByID(ctx context.Context, identity Identity, id domain.UserID) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(id),
	)
	log.Debug("Called a getByID usecase method")

	user, err := u.repo.User().GetByID(ctx, domain.UserID(id))
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	return user, nil
}
