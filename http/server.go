package http

import (
	"encoding/json"
	"net/http"

	"github.com/zen-en-tonal/mw/mail"
	"github.com/zen-en-tonal/mw/net"
)

type Registry interface {
	Update(r mail.Registry) error
	All() (*[]mail.Registry, error)
}

type State struct {
	Registry
	secret string
	host   net.Domain
}

func New(r Registry, sec string, host net.Domain) State {
	return State{r, sec, host}
}

func (s State) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.listHandler(w, r)
		return
	}
	if r.Method == "POST" {
		s.newHandler(w, r)
		return
	}
}

func (s State) isValidAuth(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer "+s.secret
}

type registry struct {
	MailAddress string `json:"mail_address"`
	Service     string `json:"service"`
}

func from(r mail.Registry, host net.Domain) registry {
	return registry{
		MailAddress: mail.NewMailAddress(r.User(), host).String(),
		Service:     r.Service(),
	}
}

func (s State) listHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isValidAuth(r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	rs, err := s.Registry.All()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var body []registry
	for _, r := range *rs {
		body = append(body, from(r, s.host))
	}
	json, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}

type req struct {
	Service string `json:"service"`
}

func (s State) newHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isValidAuth(r) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var body req
	err := json.Unmarshal(buffer, &body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d, err := net.ParseDomain(body.Service)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reg := mail.Issue(*d)
	err = s.Registry.Update(reg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(from(reg, s.host))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}