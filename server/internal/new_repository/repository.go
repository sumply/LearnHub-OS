package repository

import (
	"context"
	domain "server/internal/new_domain"
)

type (
	User interface {
		Save(context.Context, domain.UserData) (domain.UserID, error)
		Get(context.Context, *UserFilter) ([]domain.UserData, error)
		Update(context.Context, domain.UserData) error
		Delete(context.Context, domain.UserID) error
	}

	UserFilter struct {
		Limit int
	}
)
