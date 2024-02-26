package mail

import "errors"

var (
	ErrInvaildRcpt   error = errors.New("invalid rcpt")
	ErrInvaildDomain error = errors.New("invalid domain")
)
