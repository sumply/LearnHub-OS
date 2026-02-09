package domain

import (
	"errors"
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

func NewUser(firstName, lastName, email, pwdHash string, role UserRole) (User, error) {
	domainErr := NewError("user")

	if firstName == "" {
		domainErr.add("first_name", errors.New("first_name is empty"))
	}
	if lastName == "" {
		domainErr.add("last_name", errors.New("last_name is empty"))
	}
	if email == "" {
		domainErr.add("email", errors.New("email is empty"))
	}
	if pwdHash == "" {
		domainErr.add("password_hash", errors.New("password_hash is empty"))
	}
	if !role.Validate() {
		domainErr.add("role", errors.New("role is invalid"))
	}

	if !domainErr.Empty() {
		return User{}, domainErr
	}

	return User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		Email:     email,
		PwdHash:   pwdHash,
		CreatedAt: time.Now().UTC(),
	}, nil
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
	*User
	Group uuid.UUID
}

func NewStudent(u *User, groupID uuid.UUID) (Student, error) {
	if u == nil {
		panic("user is nil")
	}

	domainErr := NewError("student")
	if groupID == uuid.Nil {
		domainErr.add("group_id", errors.New("group_id is empty"))
	}

	if !domainErr.Empty() {
		return Student{}, domainErr
	}

	return Student{
		User:  u,
		Group: groupID,
	}, nil
}

type Teacher struct {
	*User
	Subjects uuid.UUIDs
	Groups   uuid.UUIDs
}

func NewTeacher(user *User, subjects, groups uuid.UUIDs) (Teacher, error) {
	if user == nil {
		panic("user is nil")
	}

	domainErr := NewError("teacher")
	if len(subjects) == 0 {
		domainErr.add("subject_ids", errors.New("subject_ids is empty"))
	}
	if len(groups) == 0 {
		domainErr.add("group_ids", errors.New("group_ids is empty"))
	}

	if !domainErr.Empty() {
		return Teacher{}, domainErr
	}

	return Teacher{
		User:     user,
		Subjects: subjects,
		Groups:   groups,
	}, nil
}
