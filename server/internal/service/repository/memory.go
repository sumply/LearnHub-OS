package repository

import "time"

type ID uint64
type Foreign ID
type UserRole string

const (
	Root    UserRole = "root"
	Admin   UserRole = "admin"
	Student UserRole = "student"
	Teacher UserRole = "teacher"
)

type Repository interface {
	User() User
}

type User interface {
	Save(UserSaveParam) error
}

type AuthEntity struct {
	Login          string
	Email          string
	PasswordHashed string
}

type ProfileEntity struct {
	ID         ID
	FirstName  string
	LastName   string
	MiddleName string
	Role       UserRole
	CreatedAt  time.Time
}

type GroupEntity struct {
	ID   ID
	Name string
}

type SubjectEntity struct {
	ID   ID
	Name string
}

type QuizEntity struct {
	ID         ID
	Title      string
	Summary    string
	TotalScore int
	CreatedAt  time.Time
}

type QuestionEntity struct {
	ID    ID
	Title string
}

type OptionEntity struct {
	ID        ID
	Text      string
	IsCorrect bool
}
type ProgressEntity struct {
	ID             ID
	Score          int
	CompletionTime time.Time
	IsCompleted    bool
}
type UserSaveParam struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Login      string
	HashedPwd  string
}

func NewStub() *Stub {
	return &Stub{}
}

type Stub struct{}

func (s *Stub) User() User {
	return &StubUser{}
}

type StubUser struct{}

func (s *StubUser) Save(UserSaveParam) error {
	return nil
}
