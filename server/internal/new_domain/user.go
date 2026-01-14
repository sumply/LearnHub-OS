package domain

import (
	"fmt"
	"regexp"
	"server/internal/common"
	"strings"
	"time"
	"unicode"
)

type User interface {
	Profile() *Profile
}

type Credential struct {
	ID      UserID
	Login   Login
	PwdHash PwdHash
	Email   Email
}

type Profile struct {
	Credential UserID
	FirstName  UserName
	LastName   UserName
	MiddleName UserName
	Role       UserRole
	CreatedAt  time.Time
}

type CoreUser struct {
	ID      UserID
	profile *Profile
}

func (u *CoreUser) Profile() *Profile {
	if u.profile == nil {
		return &Profile{}
	}
	return u.profile
}

type Teacher struct {
	ID       TeacherID
	profile  *Profile
	Subjects []SubjectID
	Groups   []GroupID
}

func (t *Teacher) Profile() *Profile {
	if t.profile == nil {
		return &Profile{}
	}
	return t.profile
}

type Student struct {
	ID      StudentID
	profile *Profile
	Group   GroupID
}

func (s *Student) Profile() *Profile {
	if s.profile == nil {
		return &Profile{}
	}
	return s.profile
}

func NewCredential(email string, maker CredentialMaker) (*Credential, Password, error) {
	login := maker.GenerateLogin()
	pwd := maker.GeneratePassword()
	hash := maker.Hash(pwd)
	_email, err := NewEmail(email)
	if err != nil {
		return nil, "", err
	}
	return &Credential{
		Login:   login,
		PwdHash: hash,
		Email:   _email,
	}, pwd, nil
}

func NewProfile(firstName, lastName, middleName string, role UserRole, credential UserID) (*Profile, error) {
	if !role.IsValid() {
		return nil, fmt.Errorf("role is invalid")
	}
	_firstName, err := newUserName(first_name, firstName)
	if err != nil {
		return nil, err
	}
	_lastName, err := newUserName(last_name, lastName)
	if err != nil {
		return nil, err
	}
	var _middleName UserName
	if middleName != "" {
		var err error
		_middleName, err = newUserName(middle_name, middleName)
		if err != nil {
			return nil, err
		}
	}
	return &Profile{
		FirstName:  _firstName,
		LastName:   _lastName,
		MiddleName: _middleName,
		Role:       role,
		Credential: credential,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func NewTeacher(user UserID, subjects []SubjectID, groups []GroupID) *Teacher {
	return &Teacher{
		ID:       TeacherID(user),
		Subjects: subjects,
		Groups:   groups,
	}
}

func NewStudent(user UserID, group GroupID) *Student {
	return &Student{
		ID:    StudentID(user),
		Group: group,
	}
}

type CredentialMaker interface {
	PwdHasher
	GenerateLogin() Login
	GeneratePassword() Password
}

type PwdHasher interface {
	Hash(Password) PwdHash
}

type UserID common.ID

type TeacherID UserID

type StudentID UserID

type PwdHash string

type Password string

type Login string

type Email string

type UserName string

func newUserName(piece UserNamePiece, name string) (UserName, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("%s is empty", piece)
	}
	runes := []rune(name)
	if len(runes) < 2 {
		return "", fmt.Errorf("%s is too short", piece)
	}
	if len(runes) > 100 {
		return "", fmt.Errorf("%s is too long", piece)
	}
	for i, r := range runes {
		if !unicode.IsLetter(r) {
			return "", fmt.Errorf("%s is not letter", piece)
		}
		runes[i] = unicode.ToLower(r)
	}
	runes[0] = unicode.ToUpper(runes[0])
	return UserName(runes), nil
}

type UserNamePiece string

const (
	first_name  UserNamePiece = "first name"
	last_name   UserNamePiece = "last name"
	middle_name UserNamePiece = "middle name"
)

var emailRegex = regexp.MustCompile(`^\S+@\S+\.\S+$`)

func NewEmail(email string) (Email, error) {
	if email == "" {
		return "", fmt.Errorf("email is empty")
	}
	if !emailRegex.MatchString(email) {
		return "", fmt.Errorf("email is invalid")
	}
	return Email(email), nil
}

type UserRole common.Enum

const (
	RoleMissing UserRole = iota
	RoleTeacher
	RoleStudent
)

func (r UserRole) IsValid() bool {
	switch r {
	case
		RoleMissing,
		RoleStudent,
		RoleTeacher:
		return true
	default:
		return false
	}
}
