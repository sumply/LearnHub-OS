package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserAggregate any

type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Role      UserRole
	Email     string
	PwdHash   string
	CreatedAt time.Time
}

func (u *User) Validate() error {
	domainErr := NewError("user")

	u.FirstName = strings.TrimSpace(u.FirstName)

	if u.FirstName == "" {
		domainErr.add("first_name", errors.New("first_name is empty"))
	}

	u.LastName = strings.TrimSpace(u.LastName)

	if u.LastName == "" {
		domainErr.add("last_name", errors.New("last_name is empty"))
	}

	u.Email = strings.TrimSpace(u.Email)

	if u.Email == "" {
		domainErr.add("email", errors.New("email is empty"))
	}

	u.PwdHash = strings.TrimSpace(u.PwdHash)

	if u.PwdHash == "" {
		domainErr.add("password_hash", errors.New("password_hash is empty"))
	}

	if !u.Role.Validate() {
		domainErr.add("role", errors.New("role is invalid"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func NewUser(firstName, lastName, email, pwdHash string, role UserRole) (User, error) {
	u := User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Email:     email,
		PwdHash:   pwdHash,
		CreatedAt: time.Now().UTC(),
	}

	if err := u.Validate(); err != nil {
		return User{}, err
	}

	return u, nil
}

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
)

func (u UserRole) Validate() bool {
	switch u {
	case RoleAdmin, RoleStudent, RoleTeacher:
		return true
	default:
		return false
	}
}

type Student struct {
	User
	Group uuid.UUID
}

func (s *Student) Validate() error {
	if s.User.Role != RoleStudent {
		panic(fmt.Errorf("student domain got user with role (%s)", s.User.Role))
	}

	domainErr := NewError("student")

	if err := s.User.Validate(); err != nil {
		domainErr.add("user", err)
	}

	if s.Group == uuid.Nil {
		domainErr.add("group_id", fmt.Errorf("group_id is empty"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func NewStudent(user User, groupID uuid.UUID) (Student, error) {
	s := Student{
		User:  user,
		Group: groupID,
	}

	if err := s.Validate(); err != nil {
		return Student{}, err
	}

	return s, nil
}

type Teacher struct {
	User
	Groups uuid.UUIDs
}

func (t *Teacher) Validate() error {
	if t.User.Role != RoleTeacher {
		panic(fmt.Errorf("teacher domain got user with role (%s)", t.User.Role))
	}

	domainErr := NewError("teacher")

	if err := t.User.Validate(); err != nil {
		domainErr.add("user", err)
	}

	if len(t.Groups) == 0 {
		domainErr.add("group_ids", errors.New("group_ids is empty"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func NewTeacher(user User, groups uuid.UUIDs) (Teacher, error) {
	t := Teacher{
		User:   user,
		Groups: groups,
	}

	if err := t.Validate(); err != nil {
		return Teacher{}, err
	}

	return t, nil
}
