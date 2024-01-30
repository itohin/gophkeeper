package mailer

import (
	"net/smtp"
)

type Mailer interface {
	SendMail(to []string, message string) error
}

type SMTPMailer struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPMailer(from, password, host, port string) *SMTPMailer {
	return &SMTPMailer{
		from:     from,
		password: password,
		host:     host,
		port:     port,
	}
}

func (m *SMTPMailer) SendMail(to []string, message string) error {
	auth := smtp.PlainAuth("", m.from, m.password, m.host)
	err := smtp.SendMail(m.host+":"+m.port, auth, m.from, to, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
