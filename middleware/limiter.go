package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bb9leko/rate-limiter/store"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var (
	rateLimitStore store.RateLimitStore
	limiters       = make(map[string]*rate.Limiter)
	limitersMu     sync.Mutex
)

// InitRateLimitStore inicializa o store usando a factory
func InitRateLimitStore() {
	var err error
	rateLimitStore, err = store.NewRateLimitStore()
	if err != nil {
		fmt.Printf("Erro ao inicializar store: %v\n", err)
		return
	}

	// Testa a conexão
	ctx := context.Background()
	if err := rateLimitStore.Ping(ctx); err != nil {
		fmt.Printf("Erro ao conectar no store: %v\n", err)
	} else {
		fmt.Println("Store inicializado e conectado com sucesso")
	}
}

func AllowRequest(key string, isToken bool) (allowed bool, status int, msg string) {
	rateLimit := getEnvInt("IP_LIMIT_RATE", 10)
	burst := getEnvInt("IP_LIMIT_BURST", 5)
	ttl := getEnvDuration("IP_LIMIT_TTL", time.Second)

	// Se for token, sobrescreve configs se existir
	if isToken {
		rateLimit, burst, ttl = getTokenConfig(key)
	}

	// Fallback para limiter em memória se store não estiver disponível
	if rateLimitStore == nil {
		limiter := getLimiterForKey(key, rateLimit, burst)
		if !limiter.Allow() {
			return false, 429, "Código HTTP: 429 Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame"
		}
		return true, 200, ""
	}

	redisKey := "ratelimit:" + key
	ctx := context.Background()
	count, _, err := rateLimitStore.Increment(ctx, redisKey, ttl)
	if err != nil {
		return false, 500, "Erro no rate limiter"
	}
	if count > burst {
		return false, 429, "Código HTTP: 429 Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame"
	}
	return true, 200, ""
}

func getLimiterForKey(key string, rateLimit, burst int) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()
	if limiter, exists := limiters[key]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(rate.Limit(rateLimit), burst)
	limiters[key] = limiter
	return limiter
}

func getTokenConfig(token string) (rate int, burst int, ttl time.Duration) {
	rate = getEnvInt("TOKEN_RATE", getEnvInt("IP_LIMIT_RATE", 10))
	burst = getEnvInt("TOKEN_BURST", getEnvInt("IP_LIMIT_BURST", 5))
	ttl = getEnvDuration("TOKEN_TTL", getEnvDuration("IP_LIMIT_TTL", time.Second))
	return
}

func getEnvInt(name string, def int) int {
	if viper.IsSet(name) {
		return viper.GetInt(name)
	}
	return def
}

func getEnvDuration(name string, def time.Duration) time.Duration {
	if viper.IsSet(name) {
		return viper.GetDuration(name)
	}
	return def
}
