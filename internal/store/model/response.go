package model

type (
	ResultResponse struct {
		Result string `json:"result"`
	}

	AllUserURLsResponse struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	BatchCreateURLsResponse struct {
		ShortURL      string `json:"short_url"`
		CorrelationID string `json:"correlation_id"`
	}
)
