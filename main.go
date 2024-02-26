package main

import (
	"flag"
	"log"
	"net/http"
	"os"

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

func init() {
	flag.StringVar(&slackUrl, "s", slackUrl, "Slack webhook url")
	flag.StringVar(&domain, "d", domain, "Domain")
	flag.StringVar(&secret, "t", secret, "Secret")
}

func main() {
	flag.Parse()

	if secret == "" {
		log.Fatal("secret must have a value")
		return
	}

	opt := badger.DefaultOptions("db")
	kv := contact.NewKV(opt)

	filter := contact.NewFilter(kv, domain)
	slack := slack.New(slackUrl)
	r := mail.NewMailbox(filter, slack, mail.NullStorage{})
	s := smtp.New(r)

	restState := h.New(kv, secret, domain)
	http.HandleFunc("/", restState.Handle)

	s.Addr = "0.0.0.0:25"
	s.Domain = domain
	s.AllowInsecureAuth = false
	s.Debug = os.Stdout
	s.ErrorLog = log.Default()

	shutdown := make(chan bool, 1)

	listenSmtp := func() {
		log.Println("Starting SMTP server at", s.Addr)
		log.Fatal(s.ListenAndServe())
		close(shutdown)
	}
	listenHttp := func() {
		log.Println("Starting HTTP server at", "0.0.0.0:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
		close(shutdown)
	}

	go listenSmtp()
	go listenHttp()
	func() {
		for _ = range shutdown {
		}
		panic("")
	}()
}
