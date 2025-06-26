package model

// URL represents the mapping between a short ID and the original long URL
// Used for storage and communication between service layers
type URL struct {
	ShortID string `json:"short_id"` // Unique identifier for the shortened URL
	LongURL string `json:"long_url"` // Original full URL to be shortened
}
