package slack

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mw/mail"
)

type Slack struct {
	url string
}

func New(url string) Slack {
	return Slack{url}
}

func payload(env enmime.Envelope) string {
	tmp := `
	{
		"text": "New message recieved.",
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "%s"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "from: %s"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "to: %s"
				}
			},
			{
				"type": "divider"
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": %#v
				}
			}
		]
	}`
	var text []rune
	for i, r := range env.Text {
		if i >= 3000 {
			break
		}
		text = append(text, r)
	}
	return fmt.Sprintf(
		tmp,
		env.GetHeader("Subject"),
		env.GetHeader("From"),
		env.GetHeader("To"),
		string(text),
	)
}

func (s Slack) Forward(e mail.Envelope) error {
	r := e.Data()
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	}

	p := payload(*env)

	resp, err := http.Post(s.url, "application/json", bytes.NewReader([]byte(p)))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("failed to send to slack")
	}

	return nil
}
