package store

import (
	"context"
	"time"
)

// RateLimitStore define a interface para qualquer mecanismo de persistência
type RateLimitStore interface {
	Increment(ctx context.Context, key string, window time.Duration) (int, time.Duration, error)
	Reset(ctx context.Context, key string) error
	Close() error
	Ping(ctx context.Context) error
}
