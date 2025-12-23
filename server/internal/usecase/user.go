package usecase

import (
	"context"
	"errors"
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

func (u *RealUser) Login(ctx context.Context, param UserLoginParam) (JWT, error) {
	log := u.loggerFromLogin(ctx, param)
	log.Debug("Called a login usecase method")

	param.trim(u.valid)
	if !param.validate(u.valid) {
		return JWT{}, ErrInvalidField
	}

	hash := u.gen.GenHashedPwd(param.Password)
	log.Debug("Generated a password hash")

	user, err := u.repo.User().GetByLoginPwd(ctx, param.Login, hash)
	if err != nil {
		return JWT{}, u.mapStorageError(err)
	}

	access, refresh := u.gen.GenJWTTokens(uint64(user.ID), string(user.Role))
	log.Debug("Generated a jwt tokens")

	return JWT{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
func (u *RealUser) Get(context.Context, Identity) ([]UserDomain, error) {
	return nil, nil
}
func (u *RealUser) Create(ctx context.Context, identity Identity, param UserCreateParam) error {
	log := u.loggerFromCreate(ctx, identity, param)

	log.Debug("Called a create usecase method")

	if !identity.isHigherOrEqual(Admin) || !identity.isHigher(param.Role) {
		log.Warn("Access permission")
		return ErrAccess
	}

	param.trim(u.valid)

	if !param.validate(u.valid) {
		log.Warn("Param fields is not valid")
		return errors.New("fields is not valid")
	}

	login := u.gen.GenLogin()
	pwd, hashed := u.gen.GenPassword()

	err := u.repo.User().Create(ctx, repository.UserCreate{
		FirstName:  param.FirstName,
		LastName:   param.LastName,
		MiddleName: param.MiddleName,
		Email:      param.Email,
		Login:      login,
		HashedPwd:  hashed,
	})
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

func (u *RealUser) GetMe(ctx context.Context, identity Identity) (UserDomain, error) {
	log := u.loggerFromGetMe(ctx, identity)
	log.Debug("Called a getMe usecase method")

	user, err := u.repo.User().GetByID(ctx, repository.ID(identity.ID))
	if err != nil {
		return UserDomain{}, u.mapStorageError(err)
	}

	return newUserDomainFromRepo(user), nil
}

func (u *RealUser) GetByID(ctx context.Context, identity Identity, id ID) (UserDomain, error) {
	log := u.loggerFromGetByID(ctx, identity, id)
	log.Debug("Called a getByID usecase method")

	user, err := u.repo.User().GetByID(ctx, repository.ID(id))
	if err != nil {
		return UserDomain{}, u.mapStorageError(err)
	}

	return newUserDomainFromRepo(user), nil
}

func (u *RealUser) Put(context.Context, Identity, ID, UserPutParam) error {
	return nil
}

func (u *RealUser) Delete(context.Context, Identity, ID) error {
	return nil
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
