package repository

import "time"

type ID uint64
type Foreign ID
type UserRole string

type AuthEntity struct {
	Login          string
	Email          string
	PasswordHashed string
}

type UserEntity struct {
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

type SpecialityEntity struct {
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
