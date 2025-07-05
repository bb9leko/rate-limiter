package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/bb9leko/rate-limiter/middleware"
	"github.com/spf13/viper"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Requisição permitida!"))
}

func initConfig() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Aviso: .env não encontrado ou não pôde ser carregado")
	}
}

func main() {
	initConfig()
	middleware.InitRateLimitStore() // Mudança aqui

	mux := http.NewServeMux()
	mux.Handle("/", middleware.RateLimitMiddleware(http.HandlerFunc(handler)))

	log.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", mux)
}
