package mail

import (
	"io"
)

type Forward interface {
	Forward(from Contact, body io.Reader) error
}

type ContactBook interface {
	Find(to MailAddress) (*Contact, error)
}

type Session struct {
	book    ContactBook
	domain  string
	forward Forward
}

func (r Session) Accept(to MailAddress) (*Accept, error) {
	cont, err := r.book.Find(to)
	if err != nil {
		return nil, err
	}
	if to != cont.AsAddress(r.domain) {
		return nil, ErrInvaildRcpt
	}
	return &Accept{*cont, r.forward}, nil
}

type Accept struct {
	from    Contact
	forward Forward
}

func (a Accept) Forward(r io.Reader) error {
	return a.forward.Forward(a.from, r)
}
