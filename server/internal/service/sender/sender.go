package sender

type Mail interface {
	Send(email, subject, body string) error
}

func NewStubMail() *StubMail {
	return &StubMail{}
}

type StubMail struct{}

func (s *StubMail) Send(email, subject, body string) error {
	return nil
}
