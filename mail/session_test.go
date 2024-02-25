package mail

import (
	"io"
	"testing"

	"github.com/zen-en-tonal/mw/net"
)

var ownDomain = net.MustParseDomain("own.com")

func TestReception(t *testing.T) {
	expectedUser := "user"
	eptectedService := net.MustParseDomain("service.com")
	s := Session{
		registries: stubRegistries(expectedUser, eptectedService),
		domain:     ownDomain,
		forward:    ForwardMock{},
	}

	fromAddr := NewMailAddress(expectedUser, eptectedService)
	_, err := s.Reception(fromAddr)
	if err != nil {
		t.Error(err)
	}
}

func TestReceptFromDefferenceSubDomain(t *testing.T) {
	expectedUser := "user"
	eptectedService := net.MustParseDomain("service.com")
	s := Session{
		registries: stubRegistries(expectedUser, eptectedService),
		domain:     ownDomain,
		forward:    ForwardMock{},
	}

	fromAddr := NewMailAddress(expectedUser, net.MustParseDomain("sub.service.com"))
	_, err := s.Reception(fromAddr)
	if err != nil {
		t.Error(err)
	}
}

func TestReceptionFailed(t *testing.T) {
	s := Session{
		registries: stubRegistries("user", net.MustParseDomain("service.com")),
		domain:     ownDomain,
		forward:    ForwardMock{},
	}

	fromAddr := NewMailAddress("user", net.MustParseDomain("invalid.com"))
	_, err := s.Reception(fromAddr)
	if err == nil {
		t.Error("should failed")
	}
}

type ForwardMock struct{}

func (f ForwardMock) Forward(from MailAddress, body io.Reader) error {
	return nil
}

func stubRegistries(user string, service net.Domain) Registries {
	return RegistryArray{
		Registry{
			service: service,
			user:    user,
		},
	}
}
