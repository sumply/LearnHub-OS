package postgres

import (
	"context"
	domain "server/internal/new_domain"
	repository "server/internal/new_repository"
)

type (
	User struct{}
)

func (u *User) Save(context.Context, *domain.User) (domain.UserID, error)
func (u *User) Get(context.Context, *repository.UserFilter) ([]domain.Profiler, error)
func (u *User) Update(context.Context, domain.Profiler) error
func (u *User) Delete(context.Context, domain.UserID) error
