package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/generator"
)

func NewRealUser(gen generator.Generator, repo *repository.Repository) *RealUser {
	return &RealUser{
		gen:  gen,
		repo: repo,
	}
}

type RealUser struct {
	usecase
	repo *repository.Repository
	gen  generator.Generator
}

func (u *RealUser) Login(ctx context.Context, param UserLoginParam) (*domain.TokenPair, error) {
	log := logger.FromCtx(ctx).
		With(
			logger.TraceFieldFromAny(param),
		)
	log.Debug("Called a login usecase method")

	login, err := domain.NewLogin(param.Login)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}
	user, err := u.repo.User().GetByLogin(ctx, login)
	if err != nil {
		log.Warn(err.Error())
		return nil, u.mapStorageError(err)
	}

	hash, err := domain.NewPwdHashed(param.Password)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}
	ok := user.Credential.Authorization(login, hash)
	if !ok {
		return nil, ErrAccess
	}

	tokens, err := domain.NewTokenPair(user)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}
	log.Debug("Generated jwt tokens")

	return tokens, nil
}

func (u *RealUser) Get(ctx context.Context, identity Identity) ([]*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase method")

	users, err := u.repo.User().GetAll(ctx)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}

	return users, nil
}

func (u *RealUser) Create(ctx context.Context, identity Identity, param UserCreateParam) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(param),
	)

	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) || !identity.Role.IsHigher(param.Role) {
		log.Warn(ErrAccess.Error())
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
		log.Warn(err.Error())
		return err
	}

	err = u.repo.User().Save(ctx, user)
	if err != nil {
		log.Warn(err.Error())
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
