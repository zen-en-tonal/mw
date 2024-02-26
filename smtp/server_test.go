package smtp

import (
	"testing"

	"github.com/zen-en-tonal/mw/mail"
)

func TestAuthShouldFails(t *testing.T) {
	mb := mail.NewMailbox(mail.NullFilter{}, mail.NullForwarder{}, mail.NullStorage{})
	s := session{"", "", mb}
	err := s.AuthPlain("", "")
	if err == nil {
		t.Errorf("auth should fails")
	}
}
