package model

import "github.com/google/uuid"

type URL struct {
	ID       uuid.UUID
	BaseURL  string
	ShortUrl string
}

func NewUrl(url string) *URL {
	u := &URL{
		ID:      uuid.New(),
		BaseURL: url,
	}
	return u
}
