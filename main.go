package main

import (
	"net/http"

	"github.com/bb9leko/rate-limiter/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Requisição permitida!"))
}

func main() {
	middleware.InitRedisStore()
	mux := http.NewServeMux()
	mux.Handle("/", middleware.RateLimitMiddleware(http.HandlerFunc(handler)))

	http.ListenAndServe(":8080", mux)
}
