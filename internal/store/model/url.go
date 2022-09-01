package model

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

var (
	ErrURLContainSpace = errors.New("url must have no spaces in it")
	ErrURLTooShort     = errors.New("url must be 4 or more chars long")
)

type URL struct {
	ID      string `json:"result,omitempty"`
	BaseURL string `json:"url"`
}

// NewURL ...
func NewURL(url string) (*URL, error) {
	u := &URL{
		BaseURL: url,
	}
	u.ShortURL()
	if err := u.Validate(); err != nil {
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
