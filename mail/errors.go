package mail

import "errors"

var (
	ErrServiceNotFound      error = errors.New("service not found")
	ErrInvaildRcpt          error = errors.New("invalid rcpt")
	ErrSubmissionNotAllowed error = errors.New("submisson not allowed")
	ErrInvaildDomain        error = errors.New("invalid domain")
	ErrInvaildProtocol      error = errors.New("invalid protocol")
)
