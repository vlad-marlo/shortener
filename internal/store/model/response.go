package model

type ResultResponse struct {
	Result string `json:"result"`
}

type AllUserURLsResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
