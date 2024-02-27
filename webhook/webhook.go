package webhook

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mw/mail"
)

type Webhook struct {
	endpoint string
	template template.Template
}

type Payload struct {
	Text    string
	HTML    string
	Subject string
	From    string
	To      string
}

func New(url string, temp string) (*Webhook, error) {
	t, err := template.New("").Parse(temp)
	if err != nil {
		return nil, err
	}
	return &Webhook{url, *t}, nil
}

func MustNew(url string, temp string) Webhook {
	w, err := New(url, temp)
	if err != nil {
		panic(err)
	}
	return *w
}

func ToPayload(e mail.Envelope) (*Payload, error) {
	env, err := enmime.ReadEnvelope(e.Data())
	if err != nil {
		return nil, err
	}
	return &Payload{
		Text:    env.Text,
		HTML:    env.HTML,
		Subject: env.GetHeader("Subject"),
		From:    env.GetHeader("From"),
		To:      env.GetHeader("To"),
	}, nil
}

func (w Webhook) Serialize(p Payload) (io.Reader, error) {
	txt := new(bytes.Buffer)
	if err := w.template.Execute(txt, p); err != nil {
		return nil, err
	}
	return txt, nil
}

func (s Webhook) Post(r io.Reader) error {
	resp, err := http.Post(s.endpoint, "application/json", r)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to post with status %s", resp.Status)
	}
	return nil
}

func (s Webhook) Forward(e mail.Envelope) error {
	payload, err := ToPayload(e)
	if err != nil {
		return err
	}
	r, err := s.Serialize(*payload)
	if err != nil {
		return err
	}
	return s.Post(r)
}
