package domain

type TokenPair struct {
	Access  string
	Refresh string
}

type GenerateTokenPair func(*User) (*TokenPair, error)

var generateTokenPair GenerateTokenPair

func NewTokenPair(user *User) (*TokenPair, error) {
	return generateTokenPair(user)
}
