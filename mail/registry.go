package mail

import (
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/zen-en-tonal/mw/net"
)

// サービスのドメイン vs ユーザー.
type Registry struct {
	service net.Domain
	user    string
}

func (r Registry) Service() string {
	return r.service.String()
}

func (r Registry) User() string {
	return r.user
}

func NewRegistry(domain net.Domain, user string) Registry {
	return Registry{
		service: domain,
		user:    user,
	}
}

// 新しいユーザーを発行する.
//
// domainをwww.example.comとした場合、
// serviceはexample.comとなる.
func Issue(domain net.Domain) Registry {
	bytes, _ := uuid.New().MarshalBinary()
	user := hex.EncodeToString(bytes)
	return NewRegistry(domain.TrimSub(), user)
}

type Registries interface {
	Find(domain net.Domain) (*Registry, error)
}

type RegistryArray []Registry

func (a RegistryArray) Find(domain net.Domain) (*Registry, error) {
	for _, addr := range a {
		if addr.service == domain {
			return &addr, nil
		}
	}
	return nil, errors.New("not found")
}
