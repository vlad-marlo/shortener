package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrURLContainSpace     = errors.New("url must have no spaces in it")
	ErrURLTooShort         = errors.New("url must be 4 or more chars long")
	ErrURLBadCorrelationID = errors.New("correlation ID must be one")
)

type URL struct {
	BaseURL string `json:"url"`
	User    string `json:"user,omitempty"`
	CorelID string `json:"-"`
	ID      string `json:"result,omitempty"`
}

// NewURL ...
func NewURL(url, user string, correlationID ...string) (*URL, error) {
	u := &URL{
		BaseURL: url,
		User:    user,
	}
	if len(correlationID) > 1 {
		return nil, ErrURLBadCorrelationID
	} else if len(correlationID) == 1 {
		u.CorelID = correlationID[0]
	}
	if err := u.ShortURL(); err != nil {
		return nil, err
	}
	return u, nil
}

// Validate ...
func (u *URL) Validate() error {
	if strings.Contains(u.BaseURL, " ") {
		return ErrURLContainSpace
	}
	if len(u.BaseURL) < 4 {
		return ErrURLTooShort
	}
	return nil
}

// ShortURL ...
func (u *URL) ShortURL() error {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	u.ID = hex.EncodeToString(b)
	if strings.ContainsAny(u.ID, "/(=)[]{}`*&^%$#@!\\") {
		return u.ShortURL()
	}
	return u.Validate()
}
