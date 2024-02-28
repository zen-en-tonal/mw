package webhook

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mw/mail"
)

type Webhook struct {
	endpoint   string
	template   template.Template
	htmlParser *func(string) string
	html       bool
}

type Payload struct {
	Text    string
	Subject string
	From    string
	To      string
}

type Option func(*Webhook)

// WithConvertMarkdown converts HTML body into markdown.
//
// If failed to converting, the original HTML body is put in.
func WithConvertMarkdown(w *Webhook) {
	converter := md.NewConverter("", true, nil)
	f := func(s string) string {
		markdown, err := converter.ConvertString(s)
		if err != nil {
			slog.Error("failed to convert html to markdown", "internal", err)
			return s
		}
		return markdown
	}
	w.htmlParser = &f
	w.html = true
}

func New(url string, temp string, options ...Option) (*Webhook, error) {
	t, err := template.New("").Parse(temp)
	if err != nil {
		return nil, err
	}
	w := Webhook{
		endpoint:   url,
		template:   *t,
		htmlParser: nil,
		html:       false,
	}
	for _, opt := range options {
		opt(&w)
	}
	return &w, nil
}

func MustNew(url string, temp string, options ...Option) Webhook {
	w, err := New(url, temp, options...)
	if err != nil {
		panic(err)
	}
	return *w
}

// ToPayload converts the Envelope into a Payload.
func (w Webhook) ToPayload(e mail.Envelope) (*Payload, error) {
	env, err := enmime.ReadEnvelope(e.Data())
	if err != nil {
		return nil, err
	}
	if env.HTML != "" && w.html {
		env.Text = env.HTML
	}
	if w.htmlParser != nil {
		f := *w.htmlParser
		env.Text = f(env.Text)
	}
	return &Payload{
		Text:    env.Text,
		Subject: env.GetHeader("Subject"),
		From:    env.GetHeader("From"),
		To:      env.GetHeader("To"),
	}, nil
}

// Serialize unmarshals the Payload into a Reader.
func (w Webhook) Serialize(p Payload) (io.Reader, error) {
	txt := new(bytes.Buffer)
	if err := w.template.Execute(txt, p); err != nil {
		return nil, err
	}
	return txt, nil
}

// Post posts the Payload as json.
//
// If respouned a error status, Post returns a error.
func (s Webhook) Post(p Payload) error {
	r, err := s.Serialize(p)
	if err != nil {
		return err
	}
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
	payload, err := s.ToPayload(e)
	if err != nil {
		return err
	}
	return s.Post(*payload)
}
