package usecase

import (
	"server/internal/domain"

	"github.com/google/uuid"
)

type Identity interface {
	ID() uuid.UUID
	Role() domain.UserRole
}
