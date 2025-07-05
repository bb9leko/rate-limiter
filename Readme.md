## Como Funciona

- Toda requisição passa primeiro pelo RateLimitMiddleware.
- Ele verifica se o limite foi excedido (por IP ou token).
- Se sim, responde com 429 e não chama o handler final.
- Se não, chama o handler final normalmente.

### Subir Redis 

```docker-compose up -d```

### Para rodar a aplição 

```go run main.go```

### Para rodar os testes 

```go test ./middleware```

### Enviar Requisições com Token no Header para a aplicação

```for i in {1..10}; do curl -H "API_KEY: token123" -X GET http://localhost:8080; echo; done```

### Enviar Requisições com IP 

```for i in {1..10}; do curl -H "X-Forwarded-For: 192.168.1.100" http://localhost:8080; echo; done```

### Multiplos IPs simulados 

```curl -H "X-Forwarded-For: 192.168.1.101" http://localhost:8080```
```curl -H "X-Forwarded-For: 192.168.1.102" http://localhost:8080```

### Enviar IPs e Token juntos 

```for i in {1..11}; do curl -H "API_KEY: token133" -H "X-Forwarded-For:127.0.0.2" -X GET http://localhost:8080; echo; done```


### Acessar container Redis

```docker exec -it meu_redis redis-cli```

### Consulta Chaves 

```KEYS *```

### Consulta Valor Atual Contador 

```GET ratelimit:token123```

### Consulta Tempo Restante de Expiração 

```TTL ratelimit:token123```

### Após Dockerizar a aplicação 

```docker compose up --build```

```cat /etc/resolv.conf | grep nameserver```   

```for i in {1..11}; do curl -H "API_KEY: token133" -H "X-Forwarded-For:127.0.0.2" -X GET http://10.255.255.254:8080; echo; done ```