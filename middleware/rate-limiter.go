package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Função auxiliar para extrair o IP real do cliente
func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx]
	}
	return ip
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		token := r.Header.Get("API_KEY")
		var key string
		var isToken bool

		if token != "" {
			key = token
			isToken = true
		} else {
			key = getClientIP(r)
			isToken = false
		}

		allowed, status, msg := AllowRequest(key, isToken)
		if !allowed {
			lrw.WriteHeader(status)
			fmt.Fprintln(lrw, msg)
			log.Printf("Key: %s | Status: %d | Header: %v", key, lrw.statusCode, r.Header)
			return
		}
		next.ServeHTTP(lrw, r)
		log.Printf("Key: %s | Status: %d | Header: %v", key, lrw.statusCode, r.Header)
	})
}
