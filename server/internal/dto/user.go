package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
}

type UserLastAttempt struct {
	User
	LastAttempt *AttemptItem `json:"last_attempt"`
}
