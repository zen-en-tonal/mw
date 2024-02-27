package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/zen-en-tonal/mw/mail"
)

func TestValidJsonFormat(t *testing.T) {
	body := "Subject: Example message\nContent-Type: text/plain\n\nhello\nworld!"
	env := mail.MustNewEnvelope("from@mail.com", "to@mail.com", strings.NewReader(body))

	slack := NewSlack("hoge")
	r, err := slack.makePayload(env)
	if err != nil {
		t.Error(err)
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	var v interface{}
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Error(err)
	}
}
