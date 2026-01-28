package domain

import (
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

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
)

type Student struct {
	User
	Group uuid.UUID
}

type Teacher struct {
	User
	Subjects uuid.UUIDs
	Groups   uuid.UUIDs
}
