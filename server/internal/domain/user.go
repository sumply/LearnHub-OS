package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DetailedUser any

type User struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Role      UserRole  `db:"role"`
	Email     string    `db:"email"`
	PwdHash   string    `db:"pwd_hash"`
	CreatedAt time.Time `db:"created_at"`
}

type UserRole int

const (
	RoleInvalid UserRole = iota
	RoleAdmin
	RoleStudent
	RoleTeacher
)

func NewUserRole(role string) (UserRole, error) {
	switch role {
	case "admin":
		return RoleAdmin, nil
	case "student":
		return RoleStudent, nil
	case "teacher":
		return RoleTeacher, nil
	default:
		return RoleInvalid, fmt.Errorf("role: %w", ErrInvalid)
	}
}

func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	case RoleStudent:
		return "student"
	case RoleTeacher:
		return "teacher"
	default:
		return "invalid"
	}
}

type Student struct {
	User
	Group uuid.UUID `db:"group_id"`
}

type Teacher struct {
	User
	Subjects uuid.UUIDs `db:"subject_ids"`
	Groups   uuid.UUIDs `db:"group_ids"`
}
