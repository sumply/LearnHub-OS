package validator

type User interface {
	Trim(string) string
	ValidName(string) bool
	ValidEmail(string) bool
	ValidPassword(string) bool
}

func NewStubUser() *StubUser {
	return &StubUser{}
}

type StubUser struct{}

func (v *StubUser) Trim(s string) string {
	return s
}
func (v *StubUser) ValidName(string) bool {
	return true
}

func (v *StubUser) ValidEmail(string) bool {
	return true
}

func (v *StubUser) ValidPassword(string) bool {
	return true
}
