FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .env
EXPOSE 8080
CMD ["./app"]