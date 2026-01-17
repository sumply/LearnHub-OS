package generator

import domain "server/internal/new_domain"

type Generator interface {
	GenPassword() string
	GenLogin() string
}

func NewStub() *Stub {
	return &Stub{}
}

type Stub struct{}

func (g *Stub) GenPassword() string {
	return "verysecret"
}

func (g *Stub) GenLogin() string {
	return "r12345"
}

type CredentialMaker struct{}

func (m *CredentialMaker) GenerateLogin() domain.Login {
	return "r12345"
}

func (m *CredentialMaker) GeneratePassword() domain.Password {
	return "verysecret"
}

func (m *CredentialMaker) Hash(pwd domain.Password) domain.PwdHash {
	return domain.PwdHash(pwd)
}

func init() {
	_ = domain.CredentialMaker(&CredentialMaker{})
}
