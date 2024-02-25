package mail

import (
	"testing"

	"github.com/zen-en-tonal/mw/net"
)

func TestParseAddress(t *testing.T) {
	actual, err := ParseMailAddress("user@example.com")
	if err != nil {
		t.Error(err)
	}

	expacted := NewMailAddress("user", net.MustParseDomain("example.com"))
	if expacted.user != actual.user || expacted.domain.String() != actual.domain.String() {
		t.Error("assertion failed", "expected", expacted, "actual", actual)
	}
}

func TestParseInvalidAddress(t *testing.T) {
	_, err := ParseMailAddress("user@example")
	if err == nil {
		t.Error("invalid")
	}
}
