package model

var id int = 1

type URL struct {
	ID       int
	BaseURL  string
	ShortUrl string
}

func NewUrl(url string) *URL {
	u := &URL{
		ID:      id,
		BaseURL: url,
	}
	id += 1
	u.shortUrl()
	return u
}
