package smtp

import (
	"errors"
	"io"

	"github.com/emersion/go-smtp"
)

func NewServer(allowed string, data func(io.Reader) error) *smtp.Server {
	return smtp.NewServer(&backend{data: data, allowed: allowed})
}

type backend struct {
	allowed string
	data    func(io.Reader) error
}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &session{data: bkd.data, allowed: bkd.allowed}, nil
}

type session struct {
	from    string
	allowed string
	data    func(io.Reader) error
}

func (s *session) AuthPlain(username, password string) error {
	return nil
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if s.allowed != to {
		return errors.New("disallowed rcpt")
	}
	return nil
}

func (s *session) Data(r io.Reader) error {
	return s.data(r)
}

func (s *session) Reset() {}

func (s *session) Logout() error {
	return nil
}
