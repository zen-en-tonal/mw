package smtp

import (
	"errors"
	"testing"
	"time"

	"github.com/zen-en-tonal/mw/mail"
)

func TestAuthShouldFails(t *testing.T) {
	mb := mail.NewMailbox(mail.NullFilter{}, mail.NullForwarder{}, mail.NullStorage{})
	s := session{"", "", mb, time.Minute}
	err := s.AuthPlain("", "")
	if err == nil {
		t.Errorf("auth should fails")
	}
}

func TestDataErrorMustMasked(t *testing.T) {
	mb := mail.NewMailbox(fallable{}, mail.NullForwarder{}, mail.NullStorage{})
	s := session{"", "", mb, time.Minute}
	err := s.Data(nil)
	if err != ErrData {
		t.Errorf("data error should be masked")
	}
}

func TestTimeout(t *testing.T) {
	mb := mail.NewMailbox(mail.NullFilter{}, tooLong{}, mail.NullStorage{})
	s := session{"", "", mb, time.Second}
	err := s.Data(nil)
	if err != ErrData {
		t.Errorf("timeout should works")
	}
}

type fallable struct{}

func (f fallable) Validate(env mail.Envelope) error {
	return errors.New("should masked")
}

type tooLong struct{}

func (f tooLong) Forward(env mail.Envelope) error {
	time.Sleep(time.Second * 2)
	return nil
}
