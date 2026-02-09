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
		domainErr.Add("first_name", errors.New("first_name is empty"))
	}
	if lastName == "" {
		domainErr.Add("last_name", errors.New("last_name is empty"))
	}
	if email == "" {
		domainErr.Add("email", errors.New("email is empty"))
	}
	if pwdHash == "" {
		domainErr.Add("password_hash", errors.New("password_hash is empty"))
	}
	if !role.Validate() {
		domainErr.Add("role", errors.New("role is invalid"))
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
	User
	Group uuid.UUID
}

type Teacher struct {
	User
	Subjects uuid.UUIDs
	Groups   uuid.UUIDs
}
