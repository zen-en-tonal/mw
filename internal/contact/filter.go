package contact

import "github.com/zen-en-tonal/mw/mail"

type ContactBook interface {
	Find(to mail.Address) (*Contact, error)
}

type Filter struct {
	ContactBook
	domain string
}

func NewFilter(book ContactBook, domain string) Filter {
	return Filter{book, domain}
}

func (f Filter) Validate(env mail.Envelope) error {
	if env.To().Domain() != f.domain {
		return mail.ErrInvaildDomain
	}
	_, err := f.Find(env.To())
	if err != nil {
		return mail.ErrInvaildRcpt
	}
	return nil
}
