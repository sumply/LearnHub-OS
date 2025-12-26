package domain

import (
	"errors"
	"strings"
)

type SubjectID uint64

type SubjectName string

func NewSubjectName(s string) (SubjectName, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", errors.New("subject name is empty")
	}
	return SubjectName(trimmed), nil
}

type Subject struct {
	ID   SubjectID
	Name SubjectName

	Speciality []*Speciality
}

func NewSubject(name string, specs []SpecialityID) (*Subject, error) {
	newName, err := NewSubjectName(name)
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, errors.New("speciality is empty")
	}
	newSpec := make([]*Speciality, len(specs))
	for i, id := range specs {
		newSpec[i] = &Speciality{
			ID: id,
		}
	}
	return &Subject{
		Name:       newName,
		Speciality: newSpec,
	}, nil
}
