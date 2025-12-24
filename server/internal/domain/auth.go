package domain

type TokenPair struct {
	Access  string
	Refresh string
}

type TokenGenerator interface {
	GenerateTokenPair(*User) (TokenPair, error)
}
