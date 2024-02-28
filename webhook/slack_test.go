package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/zen-en-tonal/mw/mail"
)

const bodyFmt string = `From: user@inbucket.org
Subject: Example message
Content-Type: text/html

%s
`

func TestValidJsonFormat(t *testing.T) {
	body := fmt.Sprintf(bodyFmt, strings.Repeat("<strong>x</strong>\n", 4000))
	env := mail.MustNewEnvelope("from@mail.com", "to@mail.com", strings.NewReader(body))

	slack := NewSlack("hoge")
	r, err := slack.MakePayload(env)
	if err != nil {
		t.Error(err)
	}

	for i := range r.Text {
		if i > 3000 {
			t.Error("text should be less than or equals 3000 lunes")
		}
	}

	j, err := slack.Serialize(*r)
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, j)
	var v interface{}
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Error(err)
	}
}
