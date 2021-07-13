package sessionlib

import (
	"context"
	"time"
)

// SessionStore represents the store strategy of session, which can be mysql, in memory or redis etc.
type SessionStore interface {
	// Get gets the value of specified key, returning the value and error if any.
	Get(ctx context.Context, key string) (string, error)
	// Set adds or updates the key-value pair with expiration specified.
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Delete deletes the key. If key doesn't exist, do nothing.
	Delete(ctx context.Context, key string) error
}
