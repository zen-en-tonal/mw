package smtp

import "errors"

var (
	ErrAuth error = errors.New("auth error")
	ErrData error = errors.New("data error")
)
