package domain

import (
	"errors"
	"strings"
	"unicode"
)

type SpecialityID uint64

type SpecialityName string

func NewSpecialityName(s string) (SpecialityName, error) {
	s = strings.TrimSpace(s)
	if len([]rune(s)) < 2 {
		return "", errors.New("too short")
	}
	if len([]rune(s)) > 100 {
		return "", errors.New("too long")
	}
	for _, r := range s {
		if !unicode.IsDigit(r) || !unicode.IsSpace(r) {
			return "", errors.New("invalid name")
		}
	}
	return SpecialityName(s), nil
}

type Speciality struct {
	ID   SpecialityID
	Name SpecialityName
}

func NewSpeciality(name string) (*Speciality, error) {
	n, err := NewSpecialityName(name)
	if err != nil {
		return nil, err
	}
	return &Speciality{Name: n}, nil
}
