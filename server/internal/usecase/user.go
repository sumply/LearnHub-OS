package usecase

import (
	"context"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/generator"
)

func NewUserReal(gen generator.Generator, repo *repository.Repository) *UserReal {
	return &UserReal{
		gen:  gen,
		repo: repo,
	}
}

type UserReal struct {
	usecase
	repo *repository.Repository
	gen  generator.Generator
}

func (u *UserReal) Login(ctx context.Context, req *dto.LoginReq) (*domain.TokenPair, error) {
	log := logger.FromCtx(ctx).
		With(
			logger.TraceFieldFromAny(req),
		)
	log.Debug("Called a login usecase method")

	login, err := domain.NewLogin(req.Login)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}
	ctx = logger.WithLoggerCtx(ctx, log)
	user, err := u.repo.User().GetByLogin(ctx, login)
	if err != nil {
		log.Warn(err.Error())
		return nil, u.mapStorageError(err)
	}

	hash := domain.HashPassword(req.Password)
	ok := user.Credential.Authorization(login, hash)
	if !ok {
		log.Warn("bad authorization")
		return nil, ErrAccess
	}

	tokens, err := domain.NewTokenPair(user)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}

	return tokens, nil
}

func (u *UserReal) Get(ctx context.Context, identity *dto.Identity) ([]*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase method")

	ctx = logger.WithLoggerCtx(ctx, log)
	users, err := u.repo.User().GetAll(ctx)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}

	return users, nil
}

func (u *UserReal) Create(ctx context.Context, identity *dto.Identity, req *dto.UserCreateReq) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(req),
	)

	log.Debug("Called a create usecase method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) || !identity.Role.IsHigher(req.Role) {
		log.Warn(ErrAccess.Error())
		return ErrAccess
	}

	login := u.gen.GenLogin()
	pwd := u.gen.GenPassword()

	user, err := domain.NewUser(
		login,
		pwd,
		req.Email,
		req.FirstName,
		req.LastName,
		*req.MiddleName,
		req.Role,
	)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	ctx = logger.WithLoggerCtx(ctx, log)
	err = u.repo.User().Save(ctx, user)
	if err != nil {
		log.Warn(err.Error())
		return u.mapStorageError(err)
	}

	return nil
}

func (u *UserReal) GetMe(ctx context.Context, identity *dto.Identity) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a getMe usecase method")

	ctx = logger.WithLoggerCtx(ctx, log)
	user, err := u.repo.User().GetByID(ctx, identity.ID)
	if err != nil {
		log.Warn(err.Error())
		return nil, u.mapStorageError(err)
	}

	return user, nil
}

func (u *UserReal) GetByID(ctx context.Context, identity *dto.Identity, id common.ID) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(id),
	)
	log.Debug("Called a getByID usecase method")

	ctx = logger.WithLoggerCtx(ctx, log)
	user, err := u.repo.User().GetByID(ctx, id)
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	return user, nil
}
