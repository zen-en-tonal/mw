package mail

import (
	"io"

	"github.com/emersion/go-smtp"
	"github.com/zen-en-tonal/mw/net"
)

func NewServer(registries Registries, domain net.Domain, forward Forward) *smtp.Server {
	return smtp.NewServer(&backend{
		forward:    forward,
		registries: registries,
		domain:     domain,
	})
}

type backend struct {
	registries Registries
	domain     net.Domain
	forward    Forward
}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &state{
		session: Session{
			registries: bkd.registries,
			domain:     bkd.domain,
			forward:    bkd.forward,
		},
	}, nil
}

type state struct {
	session   Session
	reception *Reception
	accept    *Accept
}

func (s *state) AuthPlain(username, password string) error {
	return ErrSubmissionNotAllowed
}

func (s *state) Mail(from string, opts *smtp.MailOptions) error {
	f, err := ParseMailAddress(from)
	if err != nil {
		return err
	}
	rcpt, err := s.session.Reception(*f)
	if err != nil {
		return err
	}
	s.reception = rcpt
	return nil
}

func (s *state) Rcpt(to string, opts *smtp.RcptOptions) error {
	if s.reception == nil {
		return ErrInvaildProtocol
	}
	t, err := ParseMailAddress(to)
	if err != nil {
		return err
	}
	a, err := s.reception.Accept(*t)
	if err != nil {
		return err
	}
	s.accept = a
	return nil
}

func (s *state) Data(r io.Reader) error {
	if s.accept == nil {
		return ErrInvaildProtocol
	}
	return s.accept.Forward(r)
}

func (s *state) Reset() {
	s.accept = nil
	s.reception = nil
}

func (s *state) Logout() error {
	s.Reset()
	return nil
}
