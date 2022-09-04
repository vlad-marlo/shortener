package model

type BulkCreateURLRequest struct {
	CorrelationID int64  `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
