package registries

import (
	badger "github.com/dgraph-io/badger/v4"
	"github.com/zen-en-tonal/mw/mail"
	"github.com/zen-en-tonal/mw/net"
)

type KV struct {
	opt badger.Options
}

func NewKV(opt badger.Options) KV {
	return KV{opt: opt}
}

func (k KV) Find(domain net.Domain) (*mail.Registry, error) {
	db, err := badger.Open(k.opt)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(domain.String()))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			user = string(v)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	r := mail.NewRegistry(domain, user)
	return &r, nil
}

func (k KV) Update(r mail.Registry) error {
	db, err := badger.Open(k.opt)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(r.Service()), []byte(r.User()))
	})
}

func (k KV) All() (*mail.RegistryArray, error) {
	db, err := badger.Open(k.opt)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var array []mail.Registry
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			domain := net.MustParseDomain(string(item.Key()))
			err := item.Value(func(v []byte) error {
				array = append(array, mail.NewRegistry(domain, string(v)))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return (*mail.RegistryArray)(&array), nil
}
