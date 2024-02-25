package mail

import (
	"encoding/hex"

	"github.com/google/uuid"
)

type Contact struct {
	alias string
	user  string
}

func NewContact(alias string, user string) Contact {
	return Contact{alias, user}
}

func GenearteContact(alias string) Contact {
	bytes, _ := uuid.New().MarshalBinary()
	user := hex.EncodeToString(bytes)
	return NewContact(alias, user)
}

func (c Contact) AsAddress(domain string) MailAddress {
	return NewMailAddress(c.user, domain)
}

func (c Contact) User() string {
	return c.user
}

func (c Contact) Alias() string {
	return c.alias
}
