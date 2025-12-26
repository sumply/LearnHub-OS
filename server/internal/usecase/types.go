package usecase

import (
	"context"
	"errors"
	"server/internal/domain"
)

var (
	ErrAccess       = errors.New("access permission")
	ErrNotFound     = errors.New("not found")
	ErrCollision    = errors.New("collision")
	ErrInvalidField = errors.New("invalid field")
)

type User interface {
	Login(context.Context, UserLoginParam) (*domain.TokenPair, error)
	Get(context.Context, Identity) ([]*domain.User, error)
	Create(context.Context, Identity, UserCreateParam) error
	GetMe(context.Context, Identity) (*domain.User, error)
	GetByID(context.Context, Identity, domain.UserID) (*domain.User, error)
}

type Speciality interface {
	Create(ctx context.Context, idendity Identity, name string) error
	Get(ctx context.Context, identity Identity) ([]*domain.Speciality, error)
}

type Group interface {
	Create(ctx context.Context, name string) error
	Get(ctx context.Context) ([]*domain.Group, error)
}

type Subject interface {
	Create(context.Context, Identity, SubjectCreateParam) error
	Get(ctx context.Context, auth Identity) ([]*domain.Subject, error)
}

type SubjectCreateParam struct {
	Name          string
	SpecialityIDs []domain.SpecialityID
}

type UserLoginParam struct {
	Login    string
	Password string `log:"hide"`
}

type Identity struct {
	ID   domain.UserID
	Role domain.UserRole
}

type UserCreateParam struct {
	FirstName  string
	LastName   string
	MiddleName string
	Email      string `log:"mask"`
	Role       domain.UserRole
}

/*
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

type UserLoginParam struct {
	Login    string
	Password string
}

func (p *UserLoginParam) validate(v validator.User) bool {
	return v.ValidPassword(p.Password)
}

type JWT struct {
	AccessToken  string
	RefreshToken string
}

type UserPutParam struct {
	FirstName  string
	LastName   string
	MiddleName string
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
	Role domain.UserRole
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
	User       UserDTO
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

*/
