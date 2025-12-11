package usecase

import (
	"context"
	"time"
)

type User interface {
	Login(ctx context.Context, p UserLoginParam) (JWT, error)
	Get(ctx context.Context, auth Identity) ([]UserDomain, error)
	Create(ctx context.Context, auth Identity, p UserCreateParam) error
	GetMe(ctx context.Context, auth Identity) (UserDomain, error)
	GetByID(ctx context.Context, auth Identity, id ID) (UserDomain, error)
	Put(ctx context.Context, auth Identity, id ID, p UserPutParam) error
	Delete(ctx context.Context, auth Identity, id ID) error
}

type Quiz interface {
	Create(ctx context.Context, auth Identity, param QuizCreateParam) error
	Get(ctx context.Context, auth Identity) ([]QuizDomain, error)
	GetByID(ctx context.Context, auth Identity, id ID) (QuizDomain, error)
}

type Group interface {
	Create(ctx context.Context, name string) error
	Get(ctx context.Context) ([]GroupDomain, error)
}

type Subject interface {
	Create(ctx context.Context, auth Identity, name string) error
	Get(ctx context.Context, auth Identity) ([]SubjectDomain, error)
}

type Email string

type Password string

type UserRole uint8

const (
	Root UserRole = iota
	Admin
	Teacher
	Student
)

type ID uint64

type UserLoginParam struct {
	Login    string
	Password Password
}

type JWT struct {
	AccessToken  string
	RefreshToken string
}

type UserDomain struct {
	ID         ID
	FirstName  string
	LastName   string
	MiddleName string
	Role       UserRole
	CreatedAt  time.Time
}

type UserPutParam struct {
	FirstName  string
	LastName   string
	MiddleName string
}

type UserCreateParam struct {
	FirstName  string
	LastName   string
	MiddleName string
}

type Identity struct {
	ID   ID
	Role UserRole
}

type GroupDomain struct {
	ID        ID
	Name      string
	CreatedAt time.Time
}

type SubjectDomain struct {
	ID        ID
	Name      string
	CreatedAt time.Time
}

type QuizCreateOption struct {
	Text      string
	IsCorrect bool
}

type QuizCreateQuestion struct {
	Title   string
	Options []QuizCreateOption
}

type QuizCreateParam struct {
	Name      string
	Summary   string
	SubjectID ID
	Questions []QuizCreateQuestion
}

type QuizOptionsDomain struct {
	ID        ID
	Text      string
	IsCorrect bool
}

type QuizQuestionDomain struct {
	ID      ID
	Name    string
	Answers []QuizOptionsDomain
}

type QuizDomain struct {
	ID         ID
	Name       string
	Summary    string
	Quiestions []QuizQuestionDomain
}
