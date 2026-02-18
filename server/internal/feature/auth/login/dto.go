package login

import (
	"server/internal/domain"

	"github.com/google/uuid"
)

type Input struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Output struct {
	ID   uuid.UUID       `json:"id"`
	Role domain.UserRole `json:"role"`
}
