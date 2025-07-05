package store

import (
	"fmt"
	"os"
)

// NewRateLimitStore cria uma instância do store baseado na configuração
func NewRateLimitStore() (RateLimitStore, error) {
	storeType := os.Getenv("STORE_TYPE")
	if storeType == "" {
		storeType = "redis" // padrão
	}

	switch storeType {
	case "redis":
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}
		return NewRedisStore(addr), nil
	case "memory":
		return NewMemoryStore(), nil
	default:
		return nil, fmt.Errorf("store type '%s' not supported", storeType)
	}
}
