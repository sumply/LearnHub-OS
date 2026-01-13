package domain

import (
	"errors"
	"fmt"
	"regexp"
	"server/internal/common"
	"strings"
	"time"
	"unicode"
)

type CredentialID common.ID

type Credential struct {
	ID        common.ID
	Login     Login
	PwdHashed PwdHash
	Email     Email
}

type Profile struct {
	Credential CredentialID
	FirstName  UserName
	LastName   UserName
	MiddleName UserName
	Role       UserRole
	CreatedAt  time.Time
}

type TeacherID common.ID

type Teacher struct {
	Profile
	ID      TeacherID
	Subject []SubjectID
	Groups  []GroupID
}

type StudentID common.ID

type Student struct {
	Profile
	ID    StudentID
	Group GroupID
}

func NewCredential(login, password, email string) (*Credential, error) {
	l, err := NewLogin(login)
	if err != nil {
		return nil, err
	}
	p, err := NewPwdHashed(password)
	if err != nil {
		return nil, err
	}
	e, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &Credential{
		Login:     l,
		PwdHashed: p,
		Email:     e,
	}, nil
}

func (c *Credential) Authorization(login Login, hash PwdHash) bool {
	return c.Login == login && c.PwdHashed == hash
}

type User struct {
	ID         common.ID
	FirstName  UserName
	LastName   UserName
	MiddleName UserName
	Role       UserRole
	Credential *Credential
	CreatedAt  time.Time
}

func NewUser(firstName, lastName, middleName string, role UserRole, credential *Credential) (*User, error) {
	if !role.IsValid() {
		return nil, fmt.Errorf("role (%d) is invalid", role)
	}
	f, err := NewUserName(firstName)
	if err != nil {
		return nil, fmt.Errorf("%w: first name", err)
	}
	l, err := NewUserName(lastName)
	if err != nil {
		return nil, fmt.Errorf("%w: last name", err)
	}
	var m UserName
	if middleName != "" {
		var err error
		m, err = NewUserName(middleName)
		if err != nil {
			return nil, fmt.Errorf("%w: middle name", err)
		}
	}
	return &User{
		FirstName:  f,
		LastName:   l,
		MiddleName: m,
		Role:       role,
		Credential: credential,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

type UserName string

func NewUserName(name string) (UserName, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return "", errors.New("empty")
	}

	runes := []rune(trimmed)

	if len(runes) < 2 {
		return "", errors.New("too short")
	}

	if len(runes) > 100 {
		return "", errors.New("too long")
	}

	for i, r := range runes {
		if !unicode.IsLetter(r) {
			return "", errors.New("not letter")
		}
		runes[i] = unicode.ToLower(r)
	}
	runes[0] = unicode.ToUpper(runes[0])

	return UserName(runes), nil
}

type Login string

func NewLogin(s string) (Login, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errors.New("empty")
	}
	return Login(s), nil
}

type PwdHasher func(string) string

var pwdHasher PwdHasher = func(s string) string {
	return s
}

type PwdHash string

func HashPassword(s string) PwdHash {
	return PwdHash(pwdHasher(s))
}

func NewPwdHashed(s string) (PwdHash, error) {
	s = strings.TrimSpace(s)
	if len([]rune(s)) < 8 {
		return "", errors.New("password shorter than 8")
	}
	if len([]rune(s)) > 32 {
		return "", errors.New("password longer than 32")
	}
	return PwdHash(pwdHasher(s)), nil
}

type Email string

var emailRegex = regexp.MustCompile(`^\S+@\S+\.\S+$`)

func NewEmail(s string) (Email, error) {
	if s == "" {
		return "", errors.New("email is empty")
	}
	if !emailRegex.MatchString(s) {
		return "", errors.New("invalid email")
	}
	return Email(s), nil
}

type UserRole common.Enum

const (
	UserStudent UserRole = iota
	UserTeacher
	UserAdmin
	UserRoot
)

func (r UserRole) IsValid() bool {
	switch r {
	case
		UserAdmin,
		UserRoot,
		UserStudent,
		UserTeacher:
		return true
	default:
		return false
	}
}

func (r UserRole) IsHigherOrEqual(role UserRole) bool {
	return r >= role
}

func (r UserRole) IsHigher(role UserRole) bool {
	return r > role
}

func (c Credential) Copy() *Credential {
	return &c
}

func (u User) Copy() *User {
	if u.Credential != nil {
		u.Credential = u.Credential.Copy()
	}
	return &u
}
