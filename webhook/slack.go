package webhook

import (
	"log/slog"

	"github.com/zen-en-tonal/mw/mail"
)

type Slack struct {
	Webhook
}

func NewSlack(url string) Slack {
	return Slack{MustNew(url, temp)}
}

const temp = `
{
	"text": "New message recieved.",
	"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "{{ .Subject }}"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "plain_text",
				"text": "from: {{ .From }}"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "plain_text",
				"text": "to: {{ .To }}"
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "plain_text",
				"text": "{{ .Text }}"
			}
		}
	]
}`

func trim(txt string, count int) string {
	var text []rune
	for i, r := range txt {
		if i >= count {
			slog.Info("message is trimed up to 3000 runes")
			break
		}
		text = append(text, r)
	}
	return string(text)
}

func (s Slack) Forward(e mail.Envelope) error {
	p, err := ToPayload(e)
	if err != nil {
		return err
	}
	p.Text = trim(p.Text, 3000)

	r, err := s.Serialize(*p)
	if err != nil {
		return err
	}

	return s.Post(r)
}
