package client

import (
	"fmt"
	"net/smtp"
)

type Mail interface {
	Send(email, subject, body string) error
}

type MailSMTP struct {
	from    string
	pasword string
	addr    string
	auth    smtp.Auth
}

func NewMailSMTP(from, password, host, port string) *MailSMTP {
	return &MailSMTP{
		from:    from,
		pasword: password,
		addr:    host + ":" + port,
		auth:    smtp.PlainAuth("", from, password, host),
	}
}

func (m *MailSMTP) Send(email, subject, body string) error {
	msg := fmt.Sprintf("subject: %s\n%s", subject, body)
	return smtp.SendMail(email, m.auth, m.from, []string{email}, []byte(msg))
}

func NewStubMail() *StubMail {
	return &StubMail{}
}

type StubMail struct{}

func (s *StubMail) Send(email, subject, body string) error {
	return nil
}
