package generator

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
