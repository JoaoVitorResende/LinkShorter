package domain

import (
	"net/url"
	"shortapi/internal/exception"
	"strings"
)

type Link struct {
	code        string
	originalURL string
}

func NewLink(code, rawURL string) (Link, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return Link{}, exception.ErrEmptyCode
	}

	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return Link{}, exception.ErrEmptyURL
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || parsed.Host == "" {
		return Link{}, exception.ErrInvalidURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return Link{}, exception.ErrInvalidURL
	}

	return Link{code: code, originalURL: rawURL}, nil
}

func (l Link) Code() string        { return l.code }
func (l Link) OriginalURL() string { return l.originalURL }