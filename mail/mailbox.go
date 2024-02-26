package mail

type Forwarder interface {
	Forward(env Envelope) error
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
