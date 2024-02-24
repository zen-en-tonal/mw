package main

import (
	"flag"
	"log"
	"os"

	"github.com/zen-en-tonal/mw/slack"
	"github.com/zen-en-tonal/mw/smtp"
)

var addr = "0.0.0.0:25"
var slackUrl = ""
var domain = ""
var username = ""

func init() {
	flag.StringVar(&addr, "l", addr, "Listen address")
	flag.StringVar(&slackUrl, "s", slackUrl, "Slack webhook url")
	flag.StringVar(&domain, "d", domain, "Domain")
	flag.StringVar(&username, "u", username, "Username")
}

func main() {
	flag.Parse()

	allowedRcpt := username + "@" + domain
	s := smtp.NewServer(allowedRcpt, slack.Post(slackUrl))

	s.Addr = addr
	s.Domain = domain
	s.AllowInsecureAuth = false
	s.Debug = os.Stdout

	log.Println("Starting SMTP server at", addr)
	log.Println("Allowed rcpt is", allowedRcpt)
	log.Fatal(s.ListenAndServe())
}
