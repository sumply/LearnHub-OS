package repository

import (
	"context"
	domain "server/internal/new_domain"
)

type (
	User interface {
		Save(context.Context, *domain.User) (domain.UserID, error)
		SaveStudent(context.Context, *domain.Student) (domain.StudentID, error)
		SaveTeacher(context.Context, *domain.Teacher) (domain.TeacherID, error)
		Get(context.Context, *UserFilter) ([]domain.UserData, error)
		Update(context.Context, domain.UserData) error
		Delete(context.Context, domain.UserID) error
	}

	UserFilter struct {
		Limit int
	}
)
