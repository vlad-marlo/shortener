package model

import (
	"github.com/google/uuid"
)

type URL struct {
	ID      string
	BaseURL string
}

func NewURL(url string) URL {
	return URL{
		ID:      uuid.New().String(),
		BaseURL: url,
	}
}
