package repository

import (
	"context"
	"server/internal/domain"
	"time"
)

func newStubUserEntity() UserEntity {
	return UserEntity{
		ID:         1,
		FirstName:  "FirstName",
		LastName:   "LastName",
		MiddleName: "MiddleName",
		Role:       "student",
		CreatedAt:  time.Now().UTC(),
	}
}

type stubUser struct{}

func (m *stubUser) Save(context.Context, *domain.User) error {
	return nil
}
func (m *stubUser) GetByID(ctx context.Context, id ID) (*domain.User, error) {
	user := &domain.User{
		ID:         1,
		FirstName:  "Вася",
		LastName:   "Гвоздев",
		MiddleName: "Иванович",
		Role:       domain.UserStudent,
		CreatedAt:  time.Now().UTC(),
	}
	return user, nil
}
func (m *stubUser) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	user, err := domain.NewUser(
		"login",
		"password",
		"email",
		"firstName",
		"secondName",
		"middleName",
		domain.UserAdmin,
	)
	return user, err
}

type stubStudent struct{}

func (m *stubStudent) Save(context.Context, StudentCreate) error {
	return nil
}
func (m *stubStudent) GetByGroup(context.Context, ID) ([]*domain.User, error) {
	res := make([]*domain.User, 10)
	return res, nil
}

type stubSubject struct{}

func (m *stubSubject) Save(ctx context.Context, data SubjectCreate) error {
	return nil
}

func (m *stubSubject) Get(ctx context.Context) ([]SubjectEntity, error) {
	res := make([]SubjectEntity, 5) // Примерно 5 предметов для заглушки
	for i := range res {
		res[i] = SubjectEntity{
			ID:   ID(i),
			Name: "subject",
		}
	}
	return res, nil
}

type stubGroup struct{}

func (m *stubGroup) Save(ctx context.Context, data GroupCreate) error {
	return nil
}

func (m *stubGroup) Get(context.Context) ([]GroupEntity, error) {
	res := make([]GroupEntity, 10) // Примерно 10 групп для заглушки
	for i := range res {
		res[i] = GroupEntity{
			ID:   ID(i),
			Name: "group",
		}
	}
	return res, nil
}

type stubSpeciality struct{}

func (m *stubSpeciality) Save(ctx context.Context, data *domain.Speciality) error {
	return nil
}

func (m *stubSpeciality) Get(context.Context) ([]*domain.Speciality, error) {
	res := make([]*domain.Speciality, 8) // Примерно 8 специальностей для заглушки
	return res, nil
}

type Stub struct {
}

func NewStub() *Stub {
	return &Stub{}
}

func (m *Stub) User() User {
	return new(stubUser)
}
func (m *Stub) Subject() Subject {
	return new(stubSubject)
}
func (m *Stub) Group() Group {
	return new(stubGroup)
}
func (m *Stub) Speciality() Speciality {
	return new(stubSpeciality)
}

func (m *Stub) Do(context.Context, func(Repository) error) error {
	return nil
}
