package mail

import (
	"errors"
	"regexp"
)

type MailAddress struct {
	user   string
	domain string
}

func (m MailAddress) String() string {
	return m.user + "@" + m.domain
}

func (m MailAddress) User() string {
	return m.user
}

func (m MailAddress) Domain() string {
	return m.domain
}

func NewMailAddress(user string, domain string) MailAddress {
	return MailAddress{
		user:   user,
		domain: domain,
	}
}

func ParseMailAddress(address string) (*MailAddress, error) {
	regexp := regexp.MustCompile(
		`^(?P<user>[a-zA-Z0-9_.+-]+)@(?P<domain>([a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]*\.)+[a-zA-Z]{2,})$`,
	)
	match := regexp.FindStringSubmatch(address)
	if len(match) != 4 {
		return nil, errors.New("invalid mail address")
	}
	return &MailAddress{
		user:   match[1],
		domain: match[2],
	}, nil
}
