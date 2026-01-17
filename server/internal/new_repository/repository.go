package repository

import (
	"context"
	domain "server/internal/new_domain"
)

type (
	User interface {
		Save(context.Context, domain.Profiler) (domain.UserID, error)
		Get(context.Context, *UserFilter) ([]domain.Profiler, error)
		Update(context.Context, domain.Profiler) error
		Delete(context.Context, domain.UserID) error
	}

	UserFilter struct {
		Limit int
	}
)
