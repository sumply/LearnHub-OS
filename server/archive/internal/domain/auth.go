package domain

import (
	"encoding/base64"
	"encoding/json"
)

type TokenPair struct {
	Access  string
	Refresh string
}

type GenerateTokenPair func(*User) (*TokenPair, error)

var generateTokenPair GenerateTokenPair = func(u *User) (*TokenPair, error) {
	m := map[string]any{
		"id":   u.ID,
		"role": u.Role,
	}
	token, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		Access:  base64.StdEncoding.EncodeToString(token),
		Refresh: base64.StdEncoding.EncodeToString(token),
	}, nil
}

func InitGenerateTokenPair(f GenerateTokenPair) {
	generateTokenPair = f
}

func NewTokenPair(user *User) (*TokenPair, error) {
	return generateTokenPair(user)
}
