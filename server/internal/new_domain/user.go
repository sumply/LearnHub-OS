package domain

import (
	"fmt"
	"regexp"
	"server/internal/common"
	"strings"
	"time"
	"unicode"
)

type UserData interface {
	ID() UserID
	Profile() *Profile
}

type Credential struct {
	ID      UserID
	Login   Login
	PwdHash PwdHash
	Email   Email
}

type Profile struct {
	FirstName  UserName
	LastName   UserName
	MiddleName UserName
	Role       UserRole
	Credential UserID
	Access     UserAccess
	CreatedAt  time.Time
}

type User struct {
	id      UserID
	profile *Profile
}

func (u *User) Profile() *Profile {
	if u.profile == nil {
		return &Profile{}
	}
	return u.profile
}

func (u *User) ID() UserID {
	return u.id
}

type Teacher struct {
	id       TeacherID
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

func (t *Teacher) ID() UserID {
	return UserID(t.id)
}

type Student struct {
	id      StudentID
	profile *Profile
	Group   GroupID
}

func (s *Student) Profile() *Profile {
	if s.profile == nil {
		return &Profile{}
	}
	return s.profile
}

func (s *Student) ID() UserID {
	return UserID(s.id)
}

// NewCredential create new user's credential data
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

// NewProfile create a profile data about user.
func NewProfile(firstName, lastName, middleName string, role UserRole, access UserAccess, credential UserID) (*Profile, error) {
	if !role.IsValid() {
		return nil, fmt.Errorf("role is invalid")
	}
	if !access.IsValid() {
		return nil, fmt.Errorf("access is invalid")
	}
	_firstName, err := newUserName(firstNamePiece, firstName)
	if err != nil {
		return nil, err
	}
	_lastName, err := newUserName(lastNamePiece, lastName)
	if err != nil {
		return nil, err
	}
	var _middleName UserName
	if middleName != "" {
		var err error
		_middleName, err = newUserName(middleNamePiece, middleName)
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
		Access:     access,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

// NewUser create new user. A profile pointer is mutable.
func NewUser(id UserID, profile *Profile) *User {
	if profile == nil {
		profile = new(Profile)
	}
	user := User{
		id:      id,
		profile: profile,
	}
	return &user
}

// NewTeacher create new teacher. A profile pointer is mutable.
func NewTeacher(id UserID, profile *Profile, subjects []SubjectID, groups []GroupID) *Teacher {
	if profile == nil {
		profile = new(Profile)
	}
	profile.Role = RoleTeacher
	teacher := Teacher{
		id:       TeacherID(id),
		profile:  profile,
		Subjects: subjects,
		Groups:   groups,
	}
	return &teacher
}

// NewStudent create new student. A profile pointer is mutable.
func NewStudent(id UserID, profile *Profile, group GroupID) *Student {
	if profile == nil {
		profile = new(Profile)
	}
	profile.Role = RoleStudent
	student := Student{
		id:      StudentID(id),
		profile: profile,
		Group:   group,
	}
	return &student
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

func newUserName(piece userNamePiece, name string) (UserName, error) {
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

type userNamePiece string

const (
	firstNamePiece  userNamePiece = "first name"
	lastNamePiece   userNamePiece = "last name"
	middleNamePiece userNamePiece = "middle name"
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
	RoleNone UserRole = iota
	RoleTeacher
	RoleStudent
)

func (r UserRole) IsValid() bool {
	switch r {
	case
		RoleNone,
		RoleStudent,
		RoleTeacher:
		return true
	default:
		return false
	}
}

type UserAccess common.Enum

const (
	AccessUser UserAccess = iota
	AccessAdmin
)

func (a UserAccess) IsValid() bool {
	switch a {
	case
		AccessUser,
		AccessAdmin:
		return true
	default:
		return false
	}
}
