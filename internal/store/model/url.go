package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

// vars ...
var (
	// ErrURLContainSpace ...
	ErrURLContainSpace = errors.New("url must have no spaces in it")
	// ErrURLTooShort ...
	ErrURLTooShort = errors.New("url must be 4 or more chars long")
	// ErrURLBadCorrelationID ...
	ErrURLBadCorrelationID = errors.New("correlation ID must be one")
)

// URL ...
type URL struct {
	BaseURL   string `json:"url"`
	User      string `json:"user,omitempty"`
	CorelID   string `json:"-"`
	ID        string `json:"result,omitempty"`
	IsDeleted bool   `json:"-"`
}

type URLer interface {
	GetCorrelationId() string
	GetOriginalUrl() string
}

// NewURL ...
func NewURL(url, user string, correlationID ...string) (*URL, error) {
	u := &URL{
		BaseURL:   url,
		User:      user,
		IsDeleted: false,
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
	if strings.ContainsAny(u.ID, "/(=)[]{}`*&^%$#@!\\<>|\"") {
		return u.ShortURL()
	}
	return u.Validate()
}

// GetUser ...
func (u *URL) GetUser() string {
	return u.User
}

// GetCorrelationId ...
func (u *URL) GetCorrelationId() string {
	return u.CorelID
}

// GetOriginalUrl ...
func (u *URL) GetOriginalUrl() string {
	return u.BaseURL
}
