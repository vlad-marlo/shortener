package model

import "github.com/google/uuid"

type URL struct {
	ID      uuid.UUID
	BaseURL string
}

func NewUrl(url string) URL {
	return URL{
		ID:      uuid.New(),
		BaseURL: url,
	}
}
