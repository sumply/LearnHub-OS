package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Subject struct {
	ID   uuid.UUID
	Name string
}

func NewSubject(name string) (Subject, error) {
	subject := Subject{
		ID:   uuid.New(),
		Name: name,
	}

	if err := subject.Validate(); err != nil {
		return Subject{}, err
	}

	return subject, nil
}

func (s *Subject) Validate() error {
	domainErr := NewError("subject")

	s.Name = strings.TrimSpace(s.Name)

	if s.Name == "" {
		domainErr.add("name", errors.New("name is empty"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}
