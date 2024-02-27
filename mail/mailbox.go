package mail

import (
	"sync"
)

func multi[T any](fs []func(T) error, arg T) error {
	var wg sync.WaitGroup
	wg.Add(len(fs))

	ec := make(chan error, len(fs))
	for _, f := range fs {
		go func(f func(T) error, w *sync.WaitGroup) {
			ec <- f(arg)
			w.Done()
		}(f, &wg)
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

type Forwarder interface {
	Forward(env Envelope) error
}

type Forwarders []Forwarder

func (f Forwarders) Forward(env Envelope) error {
	var fs []func(Envelope) error
	for _, x := range f {
		fs = append(fs, x.Forward)
	}
	return multi[Envelope](fs, env)
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
	if err := r.Validate(env); err != nil {
		return err
	}
	fs := []func(e Envelope) error{
		r.Forwarder.Forward,
		r.Storage.Store,
	}
	if err := multi[Envelope](fs, env); err != nil {
		return err
	}
	return nil
}
