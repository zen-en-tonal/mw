package main

import (
	"flag"
	"log"
	"os"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/zen-en-tonal/mw/forward"
	"github.com/zen-en-tonal/mw/mail"
	"github.com/zen-en-tonal/mw/net"
	"github.com/zen-en-tonal/mw/registries"
)

var addr = "0.0.0.0:25"
var slackUrl = ""
var domain = ""

func init() {
	flag.StringVar(&addr, "l", addr, "Listen address")
	flag.StringVar(&slackUrl, "s", slackUrl, "Slack webhook url")
	flag.StringVar(&domain, "d", domain, "Domain")
}

func main() {
	flag.Parse()

	opt := badger.DefaultOptions("db")
	kv := registries.NewKV(opt)
	domain := net.MustParseDomain(domain)
	slack := forward.NewSlack(slackUrl)
	s := mail.NewServer(kv, domain, slack)

	init := mail.Issue(net.MustParseDomain("example.com"))
	kv.Update(init)

	s.Addr = addr
	s.Domain = domain.String()
	s.AllowInsecureAuth = false
	s.Debug = os.Stdout

	log.Println("Starting SMTP server at", addr)
	log.Println("Accept host", init.Service())
	log.Println("Accept user", init.User())

	log.Fatal(s.ListenAndServe())
}
