package client

import (
	"bytes"
	"net"
	"net/smtp"
	"text/template"
)

type SMTP interface {
	SendCreateUserInfo(email string, page *UserCreatePage) error
}

type SMTPStub struct{}

func (SMTPStub) SendCreateUserInfo(email string, page *UserCreatePage) error {
	return nil
}

type MailMsg struct {
	From    string
	To      string
	Subject string
	Msg     string
}

type SMTPClient struct {
	from string
	auth smtp.Auth
	addr string
}

const templateContext = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

{{ .Msg }}`

const templateUserCreate = `Здравствуйте, {{ .LastName }} {{ .FirstName }} {{ .MiddleName }}!
Вы были зарегестрированы на платформе LearnHub.

Данные для входа:

Логин: {{ .Login }}
Пароль: {{ .Password }}

Это сообщение сгенерированно автоматически. На него отвечать не нужно.
`

var tmpContext *template.Template
var tmpUserCreate *template.Template

func init() {
	var err error
	tmpContext, err = template.New("context").Parse(templateContext)
	if err != nil {
		panic(err)
	}
	tmpUserCreate, err = template.New("user create").Parse(templateUserCreate)
	if err != nil {
		panic(err)
	}
}

func NewSMTPClient(from, pwd, host, port string) *SMTPClient {
	return &SMTPClient{
		from: from,
		auth: smtp.PlainAuth("", from, pwd, host),
		addr: net.JoinHostPort(host, port),
	}
}

func (c *SMTPClient) send(subject, email, msg string) error {
	tempMsg := MailMsg{
		From:    c.from,
		To:      email,
		Subject: subject,
		Msg:     msg,
	}
	var buf bytes.Buffer
	if err := tmpContext.Execute(&buf, &tempMsg); err != nil {
		return err
	}
	return smtp.SendMail(c.addr, c.auth, c.from, []string{email}, buf.Bytes())
}

func (c *SMTPClient) SendCreateUserInfo(email string, page *UserCreatePage) error {
	var buf bytes.Buffer
	if err := tmpUserCreate.Execute(&buf, page); err != nil {
		return err
	}
	return c.send("LearnHub-OS: Данные для входа", email, buf.String())
}

type UserCreatePage struct {
	FirstName  string
	LastName   string
	MiddleName string
	Login      string
	Password   string
}
