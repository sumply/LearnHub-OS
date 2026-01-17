package stub

import (
	"context"
	domain "server/internal/new_domain"
	repository "server/internal/new_repository"
	"time"
)

type (
	User struct{}
)

func (u *User) Save(context.Context, domain.Profiler) (domain.UserID, error) {
	return 1, nil
}

func (u *User) Get(context.Context, *repository.UserFilter) ([]domain.Profiler, error) {
	admin := domain.User{
		ID: 0,
	}
	admin.SetProfile(&domain.Profile{
		FirstName:  "User",
		LastName:   "User",
		MiddleName: "User",
		Role:       domain.RoleNone,
		Credential: 0,
		Access:     domain.AccessAdmin,
		CreatedAt:  time.Now().UTC(),
	})
	user := domain.User{
		ID: 1,
	}
	user.SetProfile(&domain.Profile{
		FirstName:  "User",
		LastName:   "User",
		MiddleName: "User",
		Role:       domain.RoleNone,
		Credential: 1,
		Access:     domain.AccessUser,
		CreatedAt:  time.Now().UTC(),
	})
	teacher := domain.Teacher{
		ID:       2,
		Subjects: []domain.SubjectID{1},
		Groups:   []domain.GroupID{1},
	}
	teacher.SetProfile(&domain.Profile{
		FirstName:  "Teacher",
		LastName:   "Teacher",
		MiddleName: "Teacher",
		Role:       domain.RoleTeacher,
		Access:     domain.AccessUser,
		Credential: 2,
		CreatedAt:  time.Now().UTC(),
	})
	student := domain.Student{
		ID:    3,
		Group: 1,
	}
	student.SetProfile(&domain.Profile{
		FirstName:  "Student",
		LastName:   "Student",
		MiddleName: "Student",
		Role:       domain.RoleStudent,
		Access:     domain.AccessUser,
		Credential: 2,
		CreatedAt:  time.Now().UTC(),
	})
	return []domain.Profiler{
		&admin,
		&user,
		&teacher,
		&student,
	}, nil
}

func (u *User) Update(context.Context, domain.Profiler) error {
	return nil
}

func (u *User) Delete(context.Context, domain.UserID) error {
	return nil
}
