package model

import (
	"github.com/google/uuid"
)

type URL struct {
	ID      uuid.UUID
	BaseURL string
}

func NewURL(url string) URL {
	return URL{
		ID:      uuid.New(),
		BaseURL: url,
	}
}
