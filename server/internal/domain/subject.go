package domain

import "github.com/google/uuid"

type Subject struct {
	ID   uuid.UUID
	Name string
}
