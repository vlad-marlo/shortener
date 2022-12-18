package model

// types ...
type (
	// ResultResponse ...
	ResultResponse struct {
		Result string `json:"result"`
	}

	// AllUserURLsResponse ...
	AllUserURLsResponse struct {
		OriginalURL string `json:"original_url"`
		ShortURL    string `json:"short_url"`
	}

	// BatchCreateURLsResponse ...
	BatchCreateURLsResponse struct {
		ShortURL      string `json:"short_url"`
		CorrelationID string `json:"correlation_id"`
	}
)
