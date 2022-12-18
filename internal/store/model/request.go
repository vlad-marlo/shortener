package model

// BulkCreateURLRequest ...
type BulkCreateURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
