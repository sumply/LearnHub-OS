package generator

import "encoding/json"

type Generator interface {
	GenPassword() (pwd string, hashed string)
	GenLogin() string
	GenHashedPwd(pwd string) (hashed string)
	GenJWTTokens(id uint64, role string) (access, refresh string)
}

type WelcomePageParam struct {
	FirstName  string
	LastName   string
	MiddleName string
	Login      string
	Password   string
}

func NewStubPageGenerator() *StubPageGenerator {
	return &StubPageGenerator{}
}

type PageGenerator interface {
	GenWelcomePage(WelcomePageParam) string
}

type StubPageGenerator struct{}

func (g *StubPageGenerator) GenWelcomePage(WelcomePageParam) string {
	return ""
}

func NewStub() *Stub {
	return &Stub{}
}

type Stub struct{}

func (g *Stub) GenPassword() (pwd string, hashed string) {
	return "secret", "hashed"
}

func (g *Stub) GenLogin() string {
	return "r12345"
}

func (g *Stub) GenHashedPwd(pwd string) (hashed string) {
	return "hashed"
}

func (g *Stub) GenJWTTokens(id uint64, role string) (access, refresh string) {
	token := map[string]any{
		"id":   id,
		"role": role,
	}
	d, _ := json.Marshal(&token)
	return string(d), string(d)
}
