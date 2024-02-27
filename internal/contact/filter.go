package contact

import (
	"fmt"

	"github.com/zen-en-tonal/mw/mail"
)

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
		return fmt.Errorf(
			"domain miss matched. expected domain is %s but actual is %s",
			env.To().Domain(),
			f.domain,
		)
	}
	_, err := f.Find(env.To())
	if err != nil {
		return fmt.Errorf(
			"failed to find a thing with key %s: %w",
			env.To().String(),
			err,
		)
	}
	return nil
}
