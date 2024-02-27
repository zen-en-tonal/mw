package smtp

import (
	"io"
	"log/slog"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/zen-en-tonal/mw/mail"
)

func New(m mail.Mailbox, timeout time.Duration) *smtp.Server {
	return smtp.NewServer(&backend{m, timeout})
}

type backend struct {
	mail.Mailbox
	timeout time.Duration
}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{"", "", bkd.Mailbox, bkd.timeout}, nil
}

type session struct {
	from string
	to   string
	mail.Mailbox
	timeout time.Duration
}

func (s *session) AuthPlain(username, password string) error {
	slog.Info("someone try to auth")
	return ErrAuth
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
		slog.Error("failed to create envelope", "internal", err)
		return ErrData
	}
	e := make(chan error, 1)
	go func() {
		e <- s.Recieve(*env)
		close(e)
	}()
	select {
	case err := <-e:
		if err != nil {
			slog.Error("failed to recieve process", "internal", err)
			return ErrData
		}
		return nil
	case <-time.After(s.timeout):
		slog.Error("recieve process time out")
		return ErrData
	}
}

func (s *session) Reset() {
	s.from = ""
	s.to = ""
}

func (s *session) Logout() error {
	s.Reset()
	return nil
}
