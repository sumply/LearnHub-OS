package usecase

import (
	"context"
	"errors"
	"server/internal/logger"
	"server/internal/service/generator"
	"server/internal/service/sender"
	"server/internal/service/storage"
	"server/internal/service/validator"
)

func NewRealUser(gen generator.Generator, val validator.User, ms sender.Mail, pgen generator.PageGenerator, storage storage.Storage) *RealUser {
	return &RealUser{
		gen:     gen,
		valid:   val,
		ms:      ms,
		pgen:    pgen,
		storage: storage,
	}
}

type RealUser struct {
	usecase
	gen     generator.Generator
	valid   validator.User
	ms      sender.Mail
	pgen    generator.PageGenerator
	storage storage.Storage
}

func (u *RealUser) Login(context.Context, UserLoginParam) (JWT, error) {
	return JWT{}, nil
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

	err := u.storage.User().Save(storage.UserSaveParam{
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

func (u *RealUser) GetMe(context.Context, Identity) (UserDomain, error) {
	return UserDomain{}, nil
}
func (u *RealUser) GetByID(context.Context, Identity, ID) (UserDomain, error) {
	return UserDomain{}, nil
}
func (u *RealUser) Put(context.Context, Identity, ID, UserPutParam) error {
	return nil
}
func (u *RealUser) Delete(context.Context, Identity, ID) error {
	return nil
}

func (u *RealUser) loggerFromCreate(ctx context.Context, identity Identity, param UserCreateParam) logger.Logger {
	log := logger.FromCtx(ctx)
	fields := logger.TraceField{
		Key: "Usecase",
		Value: map[string]any{
			"UserCreateParam": mapFromUserCreateParam(param),
			"Identity":        mapFromIdentity(identity),
		},
	}
	return log.With(fields)
}
