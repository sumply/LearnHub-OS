package dto

import "github.com/google/uuid"

type Group struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Curator  *User     `json:"curator"`
	Students []User    `json:"students"`
}

type Quiz struct{}
