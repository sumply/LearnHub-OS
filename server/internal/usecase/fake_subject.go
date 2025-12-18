package usecase

import (
	"context"
	"time"
)

type StubSubject struct{}

func NewStubSubject() *StubSubject {
	return &StubSubject{}
}

func (f *StubSubject) Create(ctx context.Context, _ Identity, name string) error {
	return nil
}

func (f *StubSubject) Get(ctx context.Context, _ Identity) ([]SubjectDomain, error) {
	subjects := []SubjectDomain{
		{
			ID:        1,
			Name:      "Английский язык",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Русский язык",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Информатика",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Математика",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Алгебра",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Геометрия",
			CreatedAt: time.Now(),
		},
		{
			ID:        1,
			Name:      "Обществознание",
			CreatedAt: time.Now(),
		},
	}
	return subjects, nil
}
