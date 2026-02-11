package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Group struct {
	ID   uuid.UUID
	Name string
}

func NewGroup(name string) (Group, error) {
	domainErr := NewError("group")

	if name == "" {
		domainErr.add("name", errors.New("name is empty"))
	}

	if !domainErr.Empty() {
		return Group{}, domainErr
	}

	return Group{
		ID:   uuid.New(),
		Name: name,
	}, nil
}
