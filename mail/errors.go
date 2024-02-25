package mail

import "errors"

var (
	ErrServiceNotFound      error = errors.New("")
	ErrInvaildRcpt          error = errors.New("")
	ErrSubmissionNotAllowed error = errors.New("")
	ErrInvaildDomain        error = errors.New("")
	ErrInvaildProtocol      error = errors.New("")
)
