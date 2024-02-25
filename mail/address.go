package mail

import (
	"errors"
	"regexp"

	"github.com/zen-en-tonal/mw/net"
)

type MailAddress struct {
	user   string
	domain net.Domain
}

func NewMailAddress(user string, domain net.Domain) MailAddress {
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
	u := match[1]
	d := match[2]
	domain, err := net.ParseDomain(d)
	if err != nil {
		return nil, err
	}
	return &MailAddress{
		user:   u,
		domain: *domain,
	}, err
}
