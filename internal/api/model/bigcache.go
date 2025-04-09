package model

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

type CacheCreationRequest struct {
	Key               string `json:"key"`
	Value             any    `json:"value"`
	DurationInSeconds int    `json:"duration_in_seconds"`
}

// CacheEntry represents the data that stored in BigCache
// Has two props : Value and Expiration
type CacheEntry struct {
	Value      any       `json:"value"`
	Expiration time.Time `json:"expiration"`
}

// CacheAppContext is to holds shared dependencies
type CacheAppContext struct {
	Cache *bigcache.BigCache
}

type ValidationError struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
