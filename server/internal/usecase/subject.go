package usecase

import (
	"context"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/validator"
)

type RealSubject struct {
	usecase
	val  validator.Subject
	repo repository.Repository
}

func NewRealSubject() *RealSubject {
	return &RealSubject{}
}

func (s *RealSubject) Create(ctx context.Context, identity Identity, name string) error {
	log := s.loggerFromCreate(ctx, identity, name)
	log.Debug("Called a create usecase method")

	if !identity.isHigherOrEqual(Admin) {
		return ErrAccess
	}
	name = s.val.Trim(name)
	if !s.val.ValidName(name) {
		return ErrInvalidField
	}

	err := s.repo.Subject().Save(ctx, repository.SubjectCreate{
		Name:          name,
		SpecialityIDs: []repository.ID{0},
	})
	if err != nil {
		return s.mapStorageError(err)
	}
	return nil
}

func (s *RealSubject) Get(ctx context.Context, identity Identity) ([]SubjectDomain, error) {
	log := s.loggerFromGet(ctx, identity)
	log.Debug("Called a get usecase method")

	entities, err := s.repo.Subject().Get(ctx)
	if err != nil {
		return nil, s.mapStorageError(err)
	}

	resp := make([]SubjectDomain, len(entities))
	for i, e := range entities {
		resp[i] = newSubjectDomainFromEntity(e)
	}
	return resp, nil
}

func (s *RealSubject) loggerFromCreate(ctx context.Context, identity Identity, name string) logger.Logger {
	log := logger.FromCtx(ctx)
	field := s.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
		"name":     name,
	})
	return log.With(field)
}

func (s *RealSubject) loggerFromGet(ctx context.Context, identity Identity) logger.Logger {
	log := logger.FromCtx(ctx)
	field := s.tracedFieldWithUsecase(map[string]any{
		"identity": mapFromIdentity(identity),
	})
	return log.With(field)
}
