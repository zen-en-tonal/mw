package smtp

import (
	"errors"
	"io"

	"github.com/emersion/go-smtp"
	"github.com/zen-en-tonal/mw/mail"
)

func New(m mail.Mailbox) *smtp.Server {
	return smtp.NewServer(&backend{m})
}

type backend struct {
	mail.Mailbox
}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{"", "", bkd.Mailbox}, nil
}

type session struct {
	from string
	to   string
	mail.Mailbox
}

func (s *session) AuthPlain(username, password string) error {
	return errors.New("")
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *session) Data(r io.Reader) error {
	env, err := mail.NewEnvelope(s.from, s.to, r)
	if err != nil {
		return err
	}
	return s.Recieve(*env)
}

func (s *session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *session) Logout() error {
	s.Reset()
	return nil
}
