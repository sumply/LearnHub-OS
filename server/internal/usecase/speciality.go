package usecase

import (
	"context"
	"server/internal/domain"
	"server/internal/repository"
)

type RealSpeciality struct {
	usecase
	repo repository.Repository
}

func NewRealSpeciality() *RealSpeciality {
	return &RealSpeciality{}
}

func (s *RealSpeciality) Create(ctx context.Context, identity Identity, name string) error {
	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		return ErrAccess
	}

	spec, err := domain.NewSpeciality(name)
	if err != nil {
		return err
	}

	err = s.repo.Speciality().Save(ctx, spec)
	if err != nil {
		return s.mapStorageError(err)
	}

	return nil
}

func (s *RealSpeciality) Get(ctx context.Context, identity Identity) ([]*domain.Speciality, error) {
	domains, err := s.repo.Speciality().GetAll(ctx)
	if err != nil {
		return nil, s.mapStorageError(err)
	}

	return domains, nil
}
