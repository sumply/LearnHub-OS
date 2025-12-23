package usecase

import (
	"context"
	"errors"
	"server/internal/logger"
	"server/internal/repository"
	"server/internal/service/validator"
	"time"
)

var (
	ErrAccess       = errors.New("access permission")
	ErrNotFound     = errors.New("not found")
	ErrCollision    = errors.New("collision")
	ErrInvalidField = errors.New("invalid field")
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

type Answer interface {
	Create(context.Context, Identity, AnswerCreateParam) error
	Get(context.Context, Identity) ([]AnswerDomain, error)
	GetByID(context.Context, Identity, ID) (AnswerDomain, error)
}

type Group interface {
	Create(ctx context.Context, name string) error
	Get(ctx context.Context) ([]GroupDomain, error)
}

type Subject interface {
	Create(ctx context.Context, auth Identity, name string) error
	Get(ctx context.Context, auth Identity) ([]SubjectDomain, error)
}

type Speciality interface {
	Create(ctx context.Context, name string) error
}

func NewUserRole(s string) UserRole {
	switch s {
	case "root":
		return Root
	case "admin":
		return Admin
	case "teacher":
		return Teacher
	case "student":
		return Student
	default:
		return INVALID
	}
}

type UserRole uint8

const (
	Root UserRole = iota
	Admin
	Teacher
	Student
	INVALID
)

func (u UserRole) toString() string {
	switch u {
	case Root:
		return "root"
	case Admin:
		return "admin"
	case Teacher:
		return "teacher"
	case Student:
		return "student"
	default:
		return "invalid"
	}
}

type ID uint64

type UserLoginParam struct {
	Login    string
	Password string
}

func (p *UserLoginParam) trim(v validator.User) {
	p.Login = v.Trim(p.Login)
	p.Password = v.Trim(p.Password)
}

func (p *UserLoginParam) validate(v validator.User) bool {
	return v.ValidPassword(p.Password)
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

func newUserDomainFromRepo(e repository.UserEntity) UserDomain {
	return UserDomain{
		ID:         ID(e.ID),
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		MiddleName: e.MiddleName,
		Role:       NewUserRole(string(e.Role)),
		CreatedAt:  e.CreatedAt,
	}
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
	Email      string
	Role       UserRole
}

func (p *UserCreateParam) trim(v validator.User) {
	p.FirstName = v.Trim(p.FirstName)
	p.LastName = v.Trim(p.LastName)
	p.MiddleName = v.Trim(p.MiddleName)
	p.Email = v.Trim(p.Email)
}

func (p *UserCreateParam) validate(v validator.User) bool {
	return v.ValidName(p.FirstName) && v.ValidName(p.LastName) && v.ValidName(p.MiddleName) && v.ValidEmail(p.Email)
}

type Identity struct {
	ID   ID
	Role UserRole
}

func (i *Identity) isHigherOrEqual(role UserRole) bool {
	return i.Role <= role
}

func (i *Identity) isHigher(role UserRole) bool {
	return i.Role < role
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

func newSubjectDomainFromEntity(e repository.SubjectEntity) SubjectDomain {
	return SubjectDomain{
		ID:   ID(e.ID),
		Name: e.Name,
	}
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

type AnsweredQuestion struct {
	Question QuizQuestionDomain
	Answered QuizOptionsDomain
}

type AnswerDomain struct {
	ID         ID
	TotalScore int
	Score      int
	Completed  bool
	Quiz       QuizDomain
	User       UserDomain
	Answers    []AnsweredQuestion
}

type SelectedOption struct {
	QuestionID ID
	OptionID   ID
}

type AnswerCreateParam struct {
	QuizID  ID
	Answers []SelectedOption
}

func mapFromUserCreateParam(p UserCreateParam) map[string]any {
	fields := map[string]any{
		"FirstName":  p.FirstName,
		"LastName":   p.LastName,
		"MiddleName": p.MiddleName,
		"Role":       p.Role,
	}
	logger.OnDebug(func() {
		fields["Email"] = p.Email
	})
	return fields
}

func mapFromIdentity(identity Identity) map[string]any {
	fields := map[string]any{
		"ID":   identity.ID,
		"Role": identity.Role,
	}
	return fields
}
