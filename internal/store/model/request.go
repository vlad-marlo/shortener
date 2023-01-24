package model

// BulkCreateURLRequest ...
type BulkCreateURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

func (b *BulkCreateURLRequest) GetCorrelationId() string {
	return b.CorrelationID
}

func (b *BulkCreateURLRequest) GetOriginalUrl() string {
	return b.OriginalURL
}
