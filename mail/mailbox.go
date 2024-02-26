package mail

import (
	"sync"
)

type Forwarder interface {
	Forward(env Envelope) error
}

type Forwarders []Forwarder

func (f Forwarders) Forward(env Envelope) error {
	var wg sync.WaitGroup
	wg.Add(len(f))

	ec := make(chan error, len(f))

	for _, x := range f {
		go func(x Forwarder, w *sync.WaitGroup) {
			ec <- x.Forward(env)
			w.Done()
		}(x, &wg)
	}

	wg.Wait()
	close(ec)
	for err := range ec {
		if err != nil {
			return err
		}
	}

	return nil
}

type NullForwarder struct{}

func (n NullForwarder) Forward(env Envelope) error {
	return nil
}

type Filter interface {
	Validate(env Envelope) error
}

type NullFilter struct{}

func (n NullFilter) Validate(env Envelope) error {
	return nil
}

type Storage interface {
	Store(env Envelope) error
}

type NullStorage struct{}

func (n NullStorage) Store(env Envelope) error {
	return nil
}

type Mailbox struct {
	Filter
	Forwarder
	Storage
}

func NewMailbox(filter Filter, forward Forwarder, storage Storage) Mailbox {
	return Mailbox{filter, forward, storage}
}

func (r Mailbox) Recieve(env Envelope) error {
	err := r.Validate(env)
	if err != nil {
		return err
	}
	// TODO: goroutine
	err = r.Store(env)
	if err != nil {
		return err
	}
	err = r.Forward(env)
	if err != nil {
		return err
	}
	return nil
}
