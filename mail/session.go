package mail

import (
	"io"

	"github.com/zen-en-tonal/mw/net"
)

type Forward interface {
	Forward(from MailAddress, body io.Reader) error
}

type Session struct {
	registries Registries
	domain     net.Domain
	forward    Forward
}

// Registriesにあるドメインとマッチするfromを受け入れる.
func (s Session) Reception(from MailAddress) (*Reception, error) {
	registry, err := s.registries.Find(from.domain.TrimSub())
	if err != nil {
		return nil, err
	}
	addr := NewMailAddress(registry.User(), s.domain)
	return &Reception{userAddr: addr, session: &s}, nil
}

type Reception struct {
	session  *Session
	userAddr MailAddress
}

func (r Reception) Accept(to MailAddress) (*Accept, error) {
	if to != r.userAddr {
		return nil, ErrInvaildRcpt
	}
	return &Accept{reception: &r}, nil
}

type Accept struct {
	reception *Reception
}

func (a Accept) Forward(r io.Reader) error {
	if a.reception == nil {
		return ErrInvaildProtocol
	}
	f := a.reception.session.forward
	return f.Forward(a.reception.userAddr, r)
}
