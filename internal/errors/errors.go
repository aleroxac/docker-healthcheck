package errors

import "errors"

var (
	ErrInvalidProtocol = errors.New("invalid protocol, only http or https are allowed")
	ErrInvalidHost     = errors.New("invalid host")
	ErrInvalidPort     = errors.New("invalid port; Network ports in TCP and UDP range from 0 to 65535")
	ErrInvalidPath     = errors.New("invalid path")
)
