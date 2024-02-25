package registries

import (
	badger "github.com/dgraph-io/badger/v4"
	"github.com/zen-en-tonal/mw/mail"
)

type KV struct {
	opt badger.Options
}

func NewKV(opt badger.Options) KV {
	return KV{opt: opt}
}

func (k KV) Find(addr mail.MailAddress) (*mail.Contact, error) {
	db, err := badger.Open(k.opt)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var alias string
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(addr.User()))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			alias = string(v)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	r := mail.NewContact(alias, addr.User())
	return &r, nil
}

func (k KV) Update(r mail.Contact) error {
	db, err := badger.Open(k.opt)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(r.User()), []byte(r.Alias()))
	})
}

func (k KV) All() (*[]mail.Contact, error) {
	db, err := badger.Open(k.opt)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var array []mail.Contact
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			user := string(item.Key())
			err := item.Value(func(v []byte) error {
				array = append(array, mail.NewContact(string(v), user))
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

	return &array, nil
}
