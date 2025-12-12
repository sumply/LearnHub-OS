package usecase

import (
	"context"
	"time"
)

type User interface {
	Login(context.Context, UserLoginParam) (JWT, error)
	Get(context.Context, Identity) ([]UserDomain, error)
	Create(context.Context, Identity, UserCreateParam) error
	GetMe(context.Context, Identity) (UserDomain, error)
	GetByID(context.Context, Identity, ID) (UserDomain, error)
	Put(context.Context, Identity, ID, UserPutParam) error
	Delete(context.Context, Identity, ID) error
}

type Quiz interface {
	Create(context.Context, Identity, QuizCreateParam) error
	Get(context.Context, Identity) ([]QuizDomain, error)
	GetByID(context.Context, Identity, ID) (QuizDomain, error)
}

type QuizResult interface {
	Create(context.Context, Identity, []QuizResultCreateParam) error
	Get(context.Context, Identity) error
	GetByID(context.Context, Identity, ID) error
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

type IdentityQuiz struct {
	Identity
	QuizID ID
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
	ID        ID
	Name      string
	Summary   string
	Questions []QuizQuestionDomain
}

type QuizResultShortDomain struct {
	ID         ID
	TotalScore int
	Score      int
	Completed  bool
}

type QuizResultDomain struct {
	ID         ID
	TotalScore int
	Score      int
	Completed  bool
	Quiz       QuizDomain
	User       UserDomain
}

type QuizResultCreateParam struct {
	QuestionID ID
	OptionID   ID
}
