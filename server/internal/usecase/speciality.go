package usecase

import "context"

type RealSpeciality struct {
	usecase
}

func NewRealSpeciality() *RealSpeciality {
	return &RealSpeciality{}
}

func (s *RealSpeciality) Create(ctx context.Context, name string) error {
	return nil
}
