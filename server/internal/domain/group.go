package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Group struct {
	ID   uuid.UUID
	Name string
}

func NewGroup(name string) (Group, error) {
	group := Group{
		ID:   uuid.New(),
		Name: name,
	}

	if err := group.Validate(); err != nil {
		return Group{}, err
	}

	return group, nil
}

func (g *Group) Validate() error {
	domainErr := NewError("group")

	g.Name = strings.TrimSpace(g.Name)
	if g.Name == "" {
		domainErr.add("name", errors.New("name is empty"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}
