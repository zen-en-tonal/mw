package forward

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jhillyerd/enmime"
	"github.com/zen-en-tonal/mw/mail"
)

type Slack struct {
	url string
}

func NewSlack(url string) Slack {
	return Slack{url}
}

func template(env enmime.Envelope, c mail.Contact) string {
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
					"text": "alias: %s"
				}
			},
			{
				"type": "divider"
			},
			{
				"type": "section",
				"text": {
					"type": "plain_text",
					"text": "%s"
				}
			}
		]
	}`
	return fmt.Sprintf(tmp, env.GetHeader("Subject"), env.GetHeader("From"), c.Alias(), env.Text)
}

func (s Slack) Forward(a mail.Contact, r io.Reader) error {
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.url, "application/json", bytes.NewReader([]byte(template(*env, a))))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errors.New("failed to send to slack")
	}

	return nil
}
