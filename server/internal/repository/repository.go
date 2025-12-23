package repository

import (
	"context"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type User interface {
	Create(context.Context, UserCreate) error
	GetByID(context.Context, ID) (UserEntity, error)
	GetByLoginPwd(ctx context.Context, login string, pwd string) (UserEntity, error)
}

type Student interface {
	Save(context.Context, StudentCreate) error
	GetByGroup(context.Context, ID) ([]UserEntity, error)
}

type Subject interface {
	Save(context.Context, SubjectCreate) error
	Get(context.Context) ([]SubjectEntity, error)
}

type Group interface {
	Save(context.Context, GroupCreate) error
	Get(context.Context) ([]GroupEntity, error)
}

type Speciality interface {
	Save(context.Context, SpecialityCreate) error
	Get(context.Context) ([]SpecialityEntity, error)
}

type Repository interface {
	User() User
	Student() Student
	Subject() Subject
	Group() Group
	Speciality() Speciality
}

type StudentCreate struct {
	StudentID ID
	GroupID   ID
}

type SpecialityCreate struct {
	Name string
}

type SubjectCreate struct {
	Name          string
	SpecialityIDs []ID
}

type GroupCreate struct {
	Name         string
	TeacherID    ID
	SpecialityID ID
}

type UserCreate struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Role       UserRole
	Login      string
	HashedPwd  string
}
