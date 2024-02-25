package mail

import (
	"io"

	"github.com/emersion/go-smtp"
)

func NewServer(book ContactBook, domain string, forward Forward) *smtp.Server {
	return smtp.NewServer(&backend{
		forward: forward,
		book:    book,
		domain:  domain,
	})
}

type backend struct {
	book    ContactBook
	domain  string
	forward Forward
}

func (bkd *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	s := Session{
		book:    bkd.book,
		domain:  bkd.domain,
		forward: bkd.forward,
	}
	return &state{&s, nil}, nil
}

type state struct {
	*Session
	accept *Accept
}

func (s *state) AuthPlain(username, password string) error {
	return ErrSubmissionNotAllowed
}

func (s *state) Mail(from string, opts *smtp.MailOptions) error {
	_, err := ParseMailAddress(from)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) Rcpt(to string, opts *smtp.RcptOptions) error {
	t, err := ParseMailAddress(to)
	if err != nil {
		return err
	}
	a, err := s.Accept(*t)
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
}

func (s *state) Logout() error {
	s.Reset()
	return nil
}
