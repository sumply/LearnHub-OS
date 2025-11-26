package dto

import (
	"time"
)

type UserRole string

var (
	UserRoleStudent UserRole = "student"
	UserRoleTeacher UserRole = "teacher"
	UserRoleAdmin   UserRole = "admin"
	UserRoleRoot    UserRole = "root"
)

func (u UserRole) Validate() bool {
	switch u {
	case UserRoleStudent, UserRoleTeacher, UserRoleAdmin, UserRoleRoot:
		return true
	}
	return false
}

type AuthzRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

type UserResponse struct {
	ID         int64     `json:"user_id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName string    `json:"middle_name"`
	IconRef    string    `json:"icon_ref"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type MaterialCardResponse struct {
	ID        int64     `json:"material_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Size      int       `json:"size"`
	Summary   string    `json:"summary"`
	Subject   string    `json:"subject"`
	Class     string    `json:"class"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

type AuthzResponse struct {
	JWT  JWTResponse  `json:"jwt"`
	User UserResponse `json:"user"`
}

type RegistrationRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}
