package repository

import (
	"context"
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

type stubTx struct {
}

func (tx *stubTx) Context() context.Context {
	return context.Background()
}

func (tx *stubTx) Commit() error {
	return nil
}
func (tx *stubTx) Rollback() error {
	return nil
}

type stubUser struct{}

func (m *stubUser) Create(context.Context, UserCreate) error {
	return nil
}
func (m *stubUser) GetByID(ctx context.Context, id ID) (UserEntity, error) {
	return newStubUserEntity(), nil
}
func (m *stubUser) GetByLoginPwd(ctx context.Context, login string, pwd string) (UserEntity, error) {
	return newStubUserEntity(), nil
}

type stubStudent struct{}

func (m *stubStudent) Save(context.Context, StudentCreate) error {
	return nil
}
func (m *stubStudent) GetByGroup(context.Context, ID) ([]UserEntity, error) {
	res := make([]UserEntity, 10)
	for i := range res {
		res[i] = newStubUserEntity()
	}
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

func (m *stubSpeciality) Save(ctx context.Context, data SpecialityCreate) error {
	return nil
}

func (m *stubSpeciality) Get(context.Context) ([]SpecialityEntity, error) {
	res := make([]SpecialityEntity, 8) // Примерно 8 специальностей для заглушки
	for i := range res {
		res[i] = SpecialityEntity{
			ID:   ID(i),
			Name: "speciality",
		}
	}
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
func (m *Stub) Student() Student {
	return new(stubStudent)
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
