package mail

import (
	"errors"
	"regexp"
)

type Address struct {
	user   string
	domain string
}

func (m Address) String() string {
	return m.user + "@" + m.domain
}

func (m Address) User() string {
	return m.user
}

func (m Address) Domain() string {
	return m.domain
}

func NewAddress(user string, domain string) Address {
	return Address{
		user:   user,
		domain: domain,
	}
}

func ParseAddress(address string) (*Address, error) {
	regexp := regexp.MustCompile(
		`^(?P<user>[a-zA-Z0-9_.+-]+)@(?P<domain>([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,})$`,
	)
	match := regexp.FindStringSubmatch(address)
	if len(match) != 4 {
		return nil, errors.New("invalid mail address")
	}
	return &Address{
		user:   match[1],
		domain: match[2],
	}, nil
}
