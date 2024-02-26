package mail

import (
	"errors"
	"testing"
)

func TestMailboxFail(t *testing.T) {
	m := NewMailbox(NullFilter{}, fallable{}, NullStorage{})

	err := m.Recieve(Envelope{})
	if err == nil {
		t.Errorf("should failed")
	}
}

type fallable struct{}

func (f fallable) Forward(env Envelope) error {
	return errors.New("")
}
