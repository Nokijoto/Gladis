# syntax=docker/dockerfile:1

# Etap 1. Build aplikacji Go
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY .env /app/.env
# Zależności najpierw dla cache
COPY go.mod go.sum ./
RUN go mod download

# Kopiujemy resztę źródeł
COPY . .


# Budowanie statycznej binarki z katalogu cmd/bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -buildvcs=false -ldflags="-w -s" \
    -o /out/discord-gemini-bot ./cmd/bot

# Etap 2. Minimalny runtime
FROM alpine:3.20

WORKDIR /app

# Użytkownik nieuprzywilejowany
RUN adduser -D -u 10001 appuser

# Kopiujemy binarkę
COPY --from=builder /out/discord-gemini-bot /app/discord-gemini-bot

USER appuser
ENTRYPOINT ["/app/discord-gemini-bot"]
