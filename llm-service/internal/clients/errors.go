package clients

import "errors"

var (
	ErrNotSupported    = errors.New("operation not supported by provider")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrRateLimited     = errors.New("rate limited")
	ErrContextTooLarge = errors.New("context length exceeded")
)
