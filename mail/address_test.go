package mail

import (
	"testing"
)

func TestParseAddress(t *testing.T) {
	actual, err := ParseAddress("user@example.com")
	if err != nil {
		t.Error(err)
	}

	if actual.user != "user" || actual.domain != "example.com" {
		t.Error("assertion failed")
	}
}

func TestParseInvalidAddress(t *testing.T) {
	_, err := ParseAddress("user@example")
	if err == nil {
		t.Error("invalid")
	}
}
