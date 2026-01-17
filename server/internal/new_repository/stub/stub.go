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

func (u *User) Save(context.Context, domain.UserData) (domain.UserID, error) {
	return 1, nil
}

func (u *User) Get(context.Context, *repository.UserFilter) ([]domain.UserData, error) {
	admin := domain.NewUser(1, &domain.Profile{
		FirstName:  "Admin",
		LastName:   "Admin",
		MiddleName: "Admin",
		Role:       domain.RoleNone,
		Access:     domain.AccessAdmin,
		CreatedAt:  time.Now().UTC(),
	})
	teacher := domain.NewTeacher(2, &domain.Profile{
		FirstName:  "Teacher",
		LastName:   "Teacher",
		MiddleName: "Teacher",
		Role:       domain.RoleTeacher,
		Access:     domain.AccessUser,
		CreatedAt:  time.Now().UTC(),
	}, []domain.SubjectID{1}, []domain.GroupID{1})
	student := domain.NewStudent(3, &domain.Profile{
		FirstName:  "Student",
		LastName:   "Student",
		MiddleName: "Student",
		Role:       domain.RoleStudent,
		Access:     domain.AccessUser,
		CreatedAt:  time.Now().UTC(),
	}, 1)
	return []domain.UserData{
		admin,
		teacher,
		student,
	}, nil
}

func (u *User) Update(context.Context, domain.UserData) error {
	return nil
}

func (u *User) Delete(context.Context, domain.UserID) error {
	return nil
}
