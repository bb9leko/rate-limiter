# Objetivo: 
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

Descrição:<br> O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

1. Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
    1. API_KEY: < TOKEN >
3. As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

# Requisitos:

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
    - Código HTTP: 429 
    - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

# Exemplos:

1. Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
2. Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
3. Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

# Dicas:
 - Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

# Entrega:
[X] O código-fonte completo da implementação.<br>
[ ] Documentação explicando como o rate limiter funciona e como ele pode ser configurado.<br>
[ ] Testes automatizados demonstrando a eficácia e a robustez do rate limiter.<br>
[ ] Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.<br>
[X] O servidor web deve responder na porta 8080.<br>


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