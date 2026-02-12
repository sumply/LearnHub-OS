package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Subject struct {
	ID   uuid.UUID
	Name string
}

func NewSubject(name string) (Subject, error) {
	domainErr := NewError("subject")

	if name == "" {
		domainErr.add("name", errors.New("name is empty"))
	}

	if !domainErr.Empty() {
		return Subject{}, domainErr
	}

	return Subject{
		ID:   uuid.New(),
		Name: name,
	}, nil
}
