package webhook

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/zen-en-tonal/mw/mail"
)

func TestPrepareRequest(t *testing.T) {
	body := "Subject: Example message\nContent-Type: text/plain\n\nhello"
	tmp := `{"text":"{{ .Text }}","subject":"{{ .Subject }}"}`
	w := MustNew("", tmp)
	env := mail.MustNewEnvelope("from@mail.com", "to@mail.com", strings.NewReader(body))
	payload, err := ToPayload(env)
	if err != nil {
		t.Error(err)
	}
	r, err := w.Serialize(*payload)
	if err != nil {
		t.Error(err)
	}
	if readerToString(r) != `{"text":"hello","subject":"Example message"}` {
		t.Errorf("actual %s", readerToString(r))
	}
}

func readerToString(r io.Reader) string {
	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	return buf.String()
}
