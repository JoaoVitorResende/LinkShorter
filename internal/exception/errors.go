package exception

import "errors"

var (
	ErrEmpty       = errors.New("item can't be empty")
	ErrNotFound    = errors.New("item not found")
	ErrInvalidID   = errors.New("item invalid ID")
	ErrInvalidBody = errors.New("invalid body")
	ErrInvalidURL  = errors.New("invalid url")
	ErrEmptyURL    = errors.New("url can't be empty")
	ErrEmptyCode   = errors.New("code can't be empty")
	ErrCodeTaken   = errors.New("code already in use")
)
