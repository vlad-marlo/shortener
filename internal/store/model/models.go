package model

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrURLContainSpace = errors.New("url must have no spaces in it")
)

type URL struct {
	ID      string
	BaseURL string
}

func NewURL(url string) (URL, error) {
	u := URL{
		ID:      uuid.New().String(),
		BaseURL: url,
	}

	if strings.Contains(u.BaseURL, " ") {
		return URL{}, ErrURLContainSpace
	}

	return u, nil
}
