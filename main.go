package main

import (
	"flag"
	"log/slog"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	h "github.com/zen-en-tonal/mw/http"
	"github.com/zen-en-tonal/mw/internal/contact"
	"github.com/zen-en-tonal/mw/internal/slack"
	"github.com/zen-en-tonal/mw/mail"
	"github.com/zen-en-tonal/mw/smtp"
)

var slackUrl = ""
var domain = ""
var secret = ""

var smtpHost = ""
var smtpPort = 0
var smtpUser = ""
var smtpPass = ""
var smtpTo = ""

func init() {
	flag.StringVar(&slackUrl, "s", slackUrl, "Slack webhook url")
	flag.StringVar(&domain, "d", domain, "Domain")
	flag.StringVar(&secret, "t", secret, "Secret")
}

func fowarder() mail.Forwarder {
	return mail.Forwarders([]mail.Forwarder{
		slack.New(slackUrl),
		// forward.New(smtpHost, smtpPort, smtpUser, smtpPass, smtpTo),
	})
}

func filter(kv contact.KV) mail.Filter {
	return contact.NewFilter(kv, domain)
}

func storage() mail.Storage {
	return mail.NullStorage{}
}

func main() {
	flag.Parse()

	if secret == "" {
		slog.Error("secret must have a value")
		return
	}

	opt := badger.DefaultOptions("db")
	kv := contact.NewKV(opt)

	r := mail.NewMailbox(
		filter(kv),
		fowarder(),
		storage(),
	)
	s := smtp.New(r, time.Second*5)

	restState := h.New(kv, secret, domain)
	http.HandleFunc("/", restState.Handle)

	s.Addr = "0.0.0.0:25"
	s.Domain = domain
	s.AllowInsecureAuth = false

	shutdown := make(chan error)

	listenSmtp := func() {
		slog.Info("Starting SMTP server", "addr", s.Addr)
		shutdown <- s.ListenAndServe()
		close(shutdown)
	}
	listenHttp := func() {
		slog.Info("Starting HTTP server", "addr", "0.0.0.0:8080")
		shutdown <- http.ListenAndServe(":8080", nil)
		close(shutdown)
	}

	go listenSmtp()
	go listenHttp()
	for err := range shutdown {
		slog.Error("server failed", err)
	}
}
