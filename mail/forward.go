package mail

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"
)

type MailForward struct {
	hostname string
	port     int
	username string
	password string
	to       string
}

func NewMailForward(hostname string, port int, username string, password string, to string) MailForward {
	return MailForward{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
		to:       to,
	}
}

func (m MailForward) Forward(env Envelope) error {
	from := env.From().String()
	recipients := []string{m.to}
	body := new(bytes.Buffer)
	io.Copy(body, env.Data())
	addr := fmt.Sprintf("%s:%d", m.hostname, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.hostname)
	if err := smtp.SendMail(addr, auth, from, recipients, body.Bytes()); err != nil {
		return err
	}
	return nil
}
