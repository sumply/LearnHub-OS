package validator

type Normalizer interface {
	Trim(string) string
}

type User interface {
	Normalizer
	ValidName(string) bool
	ValidEmail(string) bool
	ValidPassword(string) bool
}

type Subject interface {
	Normalizer
	ValidName(string) bool
}

func NewStubUser() *StubUser {
	return &StubUser{}
}

type StubNormalizer struct{}

func (v *StubNormalizer) Trim(s string) string {
	return s
}

type StubUser struct {
	StubNormalizer
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

type StubSubject struct {
	StubNormalizer
}

func (s *StubSubject) ValidName(string) bool {
	return true
}
