package domain

import "github.com/google/uuid"

type Group struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type GroupID uuid.UUID
