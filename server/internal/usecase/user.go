package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/generator"
	"server/internal/service/sender"
	"server/internal/service/validator"
)

func NewRealUser(gen generator.Generator, val validator.User, ms sender.Mail, pgen generator.PageGenerator, storage repository.Repository) *RealUser {
	return &RealUser{
		gen:   gen,
		valid: val,
		ms:    ms,
		pgen:  pgen,
		repo:  storage,
	}
}

type RealUser struct {
	usecase
	gen   generator.Generator
	valid validator.User
	ms    sender.Mail
	pgen  generator.PageGenerator
	repo  repository.Repository
}

func (u *RealUser) Login(ctx context.Context, param UserLoginParam) (*domain.TokenPair, error) {
	log := u.loggerFromLogin(ctx, param)
	log.Debug("Called a login usecase method")

	user, err := u.repo.User().GetByLogin(ctx, param.Login)
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	ok := user.Credential.Authorization(param.Login, param.Password)
	if !ok {
		return nil, ErrAccess
	}

	access, refresh := u.gen.GenJWTTokens(uint64(user.ID), string(user.Role))
	log.Debug("Generated a jwt tokens")

	return &domain.TokenPair{
		Access:  access,
		Refresh: refresh,
	}, nil
}
func (u *RealUser) Get(context.Context, Identity) ([]*domain.User, error) {
	return nil, nil
}
func (u *RealUser) Create(ctx context.Context, identity Identity, param UserCreateParam) error {
	log := u.loggerFromCreate(ctx, identity, param)

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

	page := u.pgen.GenWelcomePage(generator.WelcomePageParam{
		FirstName:  param.FirstName,
		LastName:   param.LastName,
		MiddleName: param.MiddleName,
		Login:      login,
		Password:   pwd,
	})
	err = u.ms.Send(param.Email, "Данные входа в аккаунт", page)
	if err != nil {
		return err
	}

	return nil
}

func (u *RealUser) GetMe(ctx context.Context, identity Identity) (*domain.User, error) {
	log := u.loggerFromGetMe(ctx, identity)
	log.Debug("Called a getMe usecase method")

	user, err := u.repo.User().GetByID(ctx, domain.UserID(identity.ID))
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	return user, nil
}

func (u *RealUser) GetByID(ctx context.Context, identity Identity, id ID) (*domain.User, error) {
	log := u.loggerFromGetByID(ctx, identity, id)
	log.Debug("Called a getByID usecase method")

	user, err := u.repo.User().GetByID(ctx, domain.UserID(id))
	if err != nil {
		return nil, u.mapStorageError(err)
	}

	return user, nil
}

func (u *RealUser) loggerFromCreate(ctx context.Context, identity Identity, param UserCreateParam) logger.Logger {
	log := logger.FromCtx(ctx)
	fields := u.tracedFieldWithUsecase(map[string]any{
		"UserCreateParam": mapFromUserCreateParam(param),
		"Identity":        mapFromIdentity(identity),
	},
	)
	return log.With(fields)
}

func (u *RealUser) loggerFromGetMe(ctx context.Context, identity Identity) logger.Logger {
	log := logger.FromCtx(ctx)
	field := u.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
	})
	return log.With(field)
}

func (u *RealUser) loggerFromLogin(ctx context.Context, param UserLoginParam) logger.Logger {
	log := logger.FromCtx(ctx)
	field := u.tracedFieldWithUsecase(map[string]any{
		"UserLoginParam": map[string]any{
			"Login":    logger.Masking(param.Login),
			"Password": logger.Masking(param.Password),
		},
	},
	)
	return log.With(field)
}

func (u *RealUser) loggerFromGetByID(ctx context.Context, identity Identity, id ID) logger.Logger {
	log := logger.FromCtx(ctx)
	field := u.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
		"id":       id,
	},
	)
	return log.With(field)
}
