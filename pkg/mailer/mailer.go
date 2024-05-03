package mailer

import (
	"net/smtp"

	"github.com/itohin/gophkeeper/pkg/logger"
)

type SMTPMailer struct {
	from     string
	password string
	host     string
	port     string
	log      logger.Logger
}

func NewSMTPMailer(from, password, host, port string, log logger.Logger) *SMTPMailer {
	return &SMTPMailer{
		from:     from,
		password: password,
		host:     host,
		port:     port,
		log:      log,
	}
}

func (m *SMTPMailer) SendMail(to []string, message string) error {
	return m.send(to, message)
}

func (m *SMTPMailer) SendMailAsync(to []string, message string) {
	go m.send(to, message)
}

func (m *SMTPMailer) send(to []string, message string) error {
	auth := smtp.PlainAuth("", m.from, m.password, m.host)
	err := smtp.SendMail(m.host+":"+m.port, auth, m.from, to, []byte(message))
	//TODO: retry send on error
	if err != nil {
		m.log.Errorf("send mail to %v error: %v", to, err)
		return err
	}
	return nil
}
