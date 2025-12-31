package domain

import (
	"errors"
	"server/internal/common"
	"strings"
)

type SubjectName string

func NewSubjectName(s string) (SubjectName, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", errors.New("subject name is empty")
	}
	return SubjectName(trimmed), nil
}

type Subject struct {
	ID   common.ID
	Name SubjectName
}

func NewSubject(name string) (*Subject, error) {
	newName, err := NewSubjectName(name)
	if err != nil {
		return nil, err
	}
	return &Subject{
		Name: newName,
	}, nil
}
