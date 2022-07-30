package model

import "github.com/google/uuid"

func (u *URL) shortUrl(length int) {
	u.ShortUrl = uuid.New()[:length]
}
