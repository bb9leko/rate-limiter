package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimitMiddleware_TokenLimit(t *testing.T) {
	limiters = make(map[string]*rate.Limiter)
	handler := RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "token123")

	rr := httptest.NewRecorder()

	// token123 permite 2 instantâneas (burst), depois 5 por segundo
	for i := 0; i < 2; i++ {
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("esperado status 200, obteve %d na tentativa %d", rr.Code, i+1)
		}
		rr = httptest.NewRecorder()
	}

	// Próximas 3 devem ser permitidas se espaçadas (rate = 5/s)
	for i := 2; i < 5; i++ {
		time.Sleep(210 * time.Millisecond) // 5 por segundo = 1 a cada 200ms
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("esperado status 200, obteve %d na tentativa %d", rr.Code, i+1)
		}
		rr = httptest.NewRecorder()
	}

	// 6ª requisição sem esperar deve ser bloqueada
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("esperado status 429, obteve %d na 6ª requisição", rr.Code)
	}
}
