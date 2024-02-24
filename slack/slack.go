package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Payload struct {
	Text string `json:"text"`
}

func Post(url string) func(io.Reader) error {
	return func(r io.Reader) error {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r)
		if err != nil {
			return err
		}

		payload, err := json.Marshal(Payload{Text: buf.String()})
		if err != nil {
			return err
		}

		resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return errors.New("")
		}

		return nil
	}
}
