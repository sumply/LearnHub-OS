package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

type UserID uint64

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

var pwdHasher PwdHasher

type PwdHashed string

func NewPwdHashed(s string) (PwdHashed, error) {
	s = strings.TrimSpace(s)
	if len([]rune(s)) < 8 {
		return "", errors.New("password shorter than 8")
	}
	if len([]rune(s)) > 32 {
		return "", errors.New("password longer than 32")
	}
	return PwdHashed(pwdHasher(s)), nil
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

type UserRole uint8

const (
	UserInvalid UserRole = iota
	UserStudent
	UserTeacher
	UserAdmin
	UserRoot
)

func (r UserRole) IsHigherOrEqual(role UserRole) bool {
	return r >= role
}

func (r UserRole) IsHigher(role UserRole) bool {
	return r > role
}

type Credential struct {
	Login     Login
	PwdHashed PwdHashed
	Email     Email
}

func (c *Credential) Authorization(login, pwd string) bool {
	hash, err := NewPwdHashed(pwd)
	if err != nil {
		return false
	}
	return string(c.Login) == login && c.PwdHashed == hash
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

type User struct {
	ID         UserID
	FirstName  UserName
	LastName   UserName
	MiddleName UserName
	Role       UserRole
	Credential *Credential
	CreatedAt  time.Time
}

func NewUser(login, pwd, email, firstName, lastName, middleName string, role UserRole) (*User, error) {
	credential, err := NewCredential(login, pwd, email)
	if err != nil {
		return nil, err
	}
	f, err := NewUserName(firstName)
	if err != nil {
		return nil, fmt.Errorf("%w: first name", err)
	}
	l, err := NewUserName(lastName)
	if err != nil {
		return nil, fmt.Errorf("%w: last name", err)
	}
	m, err := NewUserName(middleName)
	if err != nil {
		return nil, fmt.Errorf("%w: middle name", err)
	}
	return &User{
		FirstName:  f,
		LastName:   l,
		MiddleName: m,
		Role:       role,
		Credential: credential,
	}, nil
}
