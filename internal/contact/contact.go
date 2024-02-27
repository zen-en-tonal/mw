package contact

import (
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/zen-en-tonal/mw/mail"
)

type Contact struct {
	alias string
	user  string
}

func New(alias string, user string) Contact {
	return Contact{alias, user}
}

func Generate(alias string) Contact {
	bytes, _ := uuid.New().MarshalBinary()
	user := hex.EncodeToString(bytes)
	return New(alias, user)
}

func (c Contact) AsAddress(domain string) mail.Address {
	return mail.NewAddress(c.user, domain)
}

func (c Contact) User() string {
	return c.user
}

func (c Contact) Alias() string {
	return c.alias
}
