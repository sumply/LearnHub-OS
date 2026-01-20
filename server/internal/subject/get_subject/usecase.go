package get_subject

import (
	"context"
	"server/internal/domain"
)

type Postgres interface {
	Subjects(context.Context) ([]domain.Subject, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) GetSubject(ctx context.Context) (Output, error) {
	subjects, err := u.postgres.Subjects(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Subjects: make([]OutputSubject, len(subjects)),
	}
	for i := range subjects {
		output.Subjects[i] = OutputSubject{
			ID:   subjects[i].ID,
			Name: subjects[i].Name,
		}
	}

	return output, nil
}
