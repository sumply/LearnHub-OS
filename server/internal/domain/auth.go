package domain

import "encoding/json"

type TokenPair struct {
	Access  string
	Refresh string
}

type GenerateTokenPair func(*User) (*TokenPair, error)

var generateTokenPair GenerateTokenPair = func(u *User) (*TokenPair, error) {
	m := map[string]any{
		"ID":   u.ID,
		"Role": u.Role,
	}
	token, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		Access:  string(token),
		Refresh: string(token),
	}, nil
}

func NewTokenPair(user *User) (*TokenPair, error) {
	return generateTokenPair(user)
}
