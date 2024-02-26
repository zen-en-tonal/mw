package mail

import (
	"io"
)

type Envelope struct {
	from *Address
	to   *Address
	data io.Reader
}

func (e Envelope) From() Address {
	return *e.from
}

func (e Envelope) To() Address {
	return *e.to
}

func (e Envelope) Data() io.Reader {
	return e.data
}

func NewEnvelope(from string, to string, data io.Reader) (*Envelope, error) {
	f, err := ParseAddress(from)
	if err != nil {
		return nil, err
	}
	t, err := ParseAddress(to)
	if err != nil {
		return nil, err
	}
	return &Envelope{from: f, to: t, data: data}, nil
}
