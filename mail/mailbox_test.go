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

func TestForwardFails(t *testing.T) {
	fs := Forwarders([]Forwarder{
		NullForwarder{},
		fallable{},
		NullForwarder{},
	})
	err := fs.Forward(Envelope{})
	if err == nil {
		t.Errorf("should failed")
	}
}

func TestForwarders(t *testing.T) {
	fs := Forwarders([]Forwarder{
		NullForwarder{},
		NullForwarder{},
	})
	err := fs.Forward(Envelope{})
	if err != nil {
		t.Errorf("should not failed")
	}
}

type fallable struct{}

func (f fallable) Forward(env Envelope) error {
	return errors.New("")
}
